package node

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"0chain.net/core/common"
	"0chain.net/core/datastore"
	. "0chain.net/core/logging"
	"go.uber.org/zap"
)

const (
	FetchStrategyRandom  = 0
	FetchStrategyNearest = 1
)

//FetchStrategy - when fetching an entity, the strategy to use to select the peer nodes
var FetchStrategy = FetchStrategyNearest

//RequestEntity - request an entity
func (np *Pool) RequestEntity(ctx context.Context, requestor EntityRequestor, params map[string]string, handler datastore.JSONEntityReqResponderF) *Node {
	rhandler := requestor(params, handler)
	var nodes []*Node
	if FetchStrategy == FetchStrategyRandom {
		nodes = np.shuffleNodes()
	} else {
		nodes = np.GetNodesByLargeMessageTime()
	}
	for _, nd := range nodes {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if nd.Status == NodeStatusInactive {
			continue
		}
		if nd == Self.Node {
			continue
		}
		if rhandler(nd) {
			return nd
		}
	}
	return nil
}

//RequestEntityFromAll - request an entity from all the nodes
func (np *Pool) RequestEntityFromAll(ctx context.Context, requestor EntityRequestor, params map[string]string, handler datastore.JSONEntityReqResponderF) {
	rhandler := requestor(params, handler)
	var nodes []*Node
	if FetchStrategy == FetchStrategyRandom {
		nodes = np.shuffleNodes()
	} else {
		nodes = np.GetNodesByLargeMessageTime()
	}
	for _, nd := range nodes {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if nd.Status == NodeStatusInactive {
			continue
		}
		if nd == Self.Node {
			continue
		}
		rhandler(nd)
	}
}

/*SetRequestHeaders - sets the send request headers*/
func SetRequestHeaders(req *http.Request, options *SendOptions, entityMetadata datastore.EntityMetadata) bool {
	SetHeaders(req)
	if options.InitialNodeID != "" {
		req.Header.Set(HeaderInitialNodeID, options.InitialNodeID)
	}
	if entityMetadata != nil {
		req.Header.Set(HeaderRequestEntityName, entityMetadata.GetName())
	}

	if options.CODEC == 0 {
		req.Header.Set(HeaderRequestCODEC, CodecJSON)
	} else {
		req.Header.Set(HeaderRequestCODEC, CodecMsgpack)
	}
	return true
}

//RequestEntityHandler - a handler that requests an entity and uses it
func RequestEntityHandler(uri string, options *SendOptions, entityMetadata datastore.EntityMetadata) EntityRequestor {
	return func(params map[string]string, handler datastore.JSONEntityReqResponderF) SendHandler {
		return func(provider *Node) bool {
			timer := provider.GetTimer(uri)
			timeout := 500 * time.Millisecond
			if options.Timeout > 0 {
				timeout = options.Timeout
			}
			url := provider.GetN2NURLBase() + uri
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return false
			}
			q := req.URL.Query()
			for k, v := range params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			if options.Compress {
				req.Header.Set("Content-Encoding", compDecomp.Encoding())
			}
			delay := common.InduceDelay()
			eName := ""
			if entityMetadata != nil {
				eName = entityMetadata.GetName()
			}
			N2n.Debug("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.String("handler", uri), zap.String("entity", eName), zap.Any("params", params), zap.Any("delay", delay))
			SetRequestHeaders(req, options, entityMetadata)
			ctx, cancel := context.WithCancel(context.TODO())
			req = req.WithContext(ctx)
			// Keep the number of messages to a node bounded
			provider.Grab()
			time.AfterFunc(timeout, cancel)
			ts := time.Now()
			Self.Node.LastActiveTime = ts
			resp, err := httpClient.Do(req)
			provider.Release()
			duration := time.Since(ts)

			if err != nil {
				N2n.Error("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri), zap.String("entity", eName), zap.Any("params", params), zap.Error(err))
				return false
			}
			if resp.StatusCode != http.StatusOK {
				readAndClose(resp.Body)
				N2n.Error("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri), zap.String("entity", eName), zap.Any("params", params), zap.Any("status_code", resp.StatusCode))
				return false
			}
			if entityMetadata == nil {
				eName = resp.Header.Get(HeaderRequestEntityName)
				if eName == "" {
					N2n.Error("requesting - no entity name in header", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri))
				}
				entityMetadata = datastore.GetEntityMetadata(eName)
				if entityMetadata == nil {
					readAndClose(resp.Body)
					N2n.Error("requesting - unknown entity", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri), zap.String("entity", eName))
					return false
				}
			}
			provider.Status = NodeStatusActive
			provider.LastActiveTime = time.Now()
			entity, err := getResponseEntity(resp, entityMetadata)
			if err != nil {
				N2n.Error("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri), zap.String("entity", eName), zap.Any("params", params), zap.Error(err))
				return false
			}
			duration = time.Since(ts)
			timer.UpdateSince(ts)
			N2n.Info("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", duration), zap.String("handler", uri), zap.String("entity", eName), zap.Any("id", entity.GetKey()), zap.Any("params", params), zap.String("codec", resp.Header.Get(HeaderRequestCODEC)))
			if delay > 0 {
				N2n.Debug("response received", zap.Int("from", provider.SetIndex), zap.Int("to", Self.SetIndex), zap.String("handler", uri), zap.String("entity", eName), zap.Any("params", params), zap.Any("delay", delay))
			}
			ctx = context.TODO()
			_, err = handler(ctx, entity)
			if err != nil {
				N2n.Error("requesting", zap.Int("from", Self.SetIndex), zap.Int("to", provider.SetIndex), zap.Duration("duration", time.Since(ts)), zap.String("handler", uri), zap.String("entity", entityMetadata.GetName()), zap.Any("params", params), zap.Error(err))
				return false
			}
			return true
		}
	}
}

func validateRequest(sender *Node, r *http.Request) bool {
	if !validateChain(sender, r) {
		return false
	}
	if !validateEntityMetadata(sender, r) {
		return false
	}
	return true
}

/*ToN2NSendEntityHandler - takes a handler that accepts an entity, processes and responds and converts it
* into something suitable for Node 2 Node communication*/
func ToN2NSendEntityHandler(handler common.JSONResponderF) common.ReqRespHandlerf {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeID := r.Header.Get(HeaderNodeID)
		sender := GetNode(nodeID)
		if sender == nil {
			N2n.Error("message received - request from unrecognized node", zap.String("from", nodeID), zap.Int("to", Self.SetIndex), zap.String("handler", r.RequestURI))
			return
		}
		if !validateRequest(sender, r) {
			return
		}
		sender.Received++
		ctx := context.TODO()
		ts := time.Now()
		data, err := handler(ctx, r)
		if err != nil {
			common.Respond(w, nil, err)
			N2n.Error("message received", zap.Int("from", sender.SetIndex), zap.Int("to", Self.SetIndex), zap.String("handler", r.RequestURI), zap.Error(err))
			return
		}
		options := &SendOptions{Compress: true}
		var buffer *bytes.Buffer
		uri := r.URL.Path
		switch v := data.(type) {
		case datastore.Entity:
			entity := v
			codec := r.Header.Get(HeaderRequestCODEC)
			switch codec {
			case "JSON":
				options.CODEC = CODEC_JSON
			case "Msgpack":
				options.CODEC = CODEC_MSGPACK
			}
			w.Header().Set(HeaderRequestCODEC, codec)
			buffer = getResponseData(options, entity)
		case *pushDataCacheEntry:
			options.CODEC = v.Options.CODEC
			if options.CODEC == 0 {
				w.Header().Set(HeaderRequestCODEC, CodecJSON)
			} else {
				w.Header().Set(HeaderRequestCODEC, CodecMsgpack)
			}
			w.Header().Set(HeaderRequestEntityName, v.EntityName)
			buffer = bytes.NewBuffer(v.Data)
			uri = r.FormValue("_puri")
		}
		if options.Compress {
			w.Header().Set("Content-Encoding", compDecomp.Encoding())
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			if cerr, ok := err.(*common.Error); ok {
				w.Header().Set(common.AppErrorHeader, cerr.Code)
			}
			http.Error(w, err.Error(), 400)
			return
		}
		sdata := buffer.Bytes()
		w.Write(sdata)
		if isPullRequest(r) {
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			updatePullStats(sender, uri, len(sdata), ts)
		}
		N2n.Info("message received", zap.Int("from", sender.SetIndex), zap.Int("to", Self.SetIndex), zap.String("handler", r.RequestURI), zap.Duration("duration", time.Since(ts)), zap.Int("codec", options.CODEC))
	}
}