package counter

import (
	"context"
	counter "github.com/atomix/api/go/atomix/primitive/counter"
	"github.com/atomix/go-framework/pkg/atomix/proxy/passthrough"
	"github.com/atomix/go-framework/pkg/atomix/util/logging"
	"google.golang.org/grpc"
)

const Type = "Counter"

const (
	setOp         = "Set"
	getOp         = "Get"
	incrementOp   = "Increment"
	decrementOp   = "Decrement"
	checkAndSetOp = "CheckAndSet"
	snapshotOp    = "Snapshot"
	restoreOp     = "Restore"
)

// RegisterProxy registers the primitive on the given node
func RegisterProxy(node *passthrough.Node) {
	node.RegisterProxy("Counter", func(server *grpc.Server, client *passthrough.Client) {
		counter.RegisterCounterServiceServer(server, &Proxy{
			Proxy: passthrough.NewProxy(client),
			log:   logging.GetLogger("atomix", "counter"),
		})
	})
}

type Proxy struct {
	*passthrough.Proxy
	log logging.Logger
}

func (s *Proxy) Set(ctx context.Context, request *counter.SetRequest) (*counter.SetResponse, error) {
	s.log.Debugf("Received SetRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Set(ctx, request)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending SetResponse %+v", response)
	return response, nil
}

func (s *Proxy) Get(ctx context.Context, request *counter.GetRequest) (*counter.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Get(ctx, request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *Proxy) Increment(ctx context.Context, request *counter.IncrementRequest) (*counter.IncrementResponse, error) {
	s.log.Debugf("Received IncrementRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Increment(ctx, request)
	if err != nil {
		s.log.Errorf("Request IncrementRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending IncrementResponse %+v", response)
	return response, nil
}

func (s *Proxy) Decrement(ctx context.Context, request *counter.DecrementRequest) (*counter.DecrementResponse, error) {
	s.log.Debugf("Received DecrementRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Decrement(ctx, request)
	if err != nil {
		s.log.Errorf("Request DecrementRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending DecrementResponse %+v", response)
	return response, nil
}

func (s *Proxy) CheckAndSet(ctx context.Context, request *counter.CheckAndSetRequest) (*counter.CheckAndSetResponse, error) {
	s.log.Debugf("Received CheckAndSetRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.CheckAndSet(ctx, request)
	if err != nil {
		s.log.Errorf("Request CheckAndSetRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending CheckAndSetResponse %+v", response)
	return response, nil
}

func (s *Proxy) Snapshot(ctx context.Context, request *counter.SnapshotRequest) (*counter.SnapshotResponse, error) {
	s.log.Debugf("Received SnapshotRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Snapshot(ctx, request)
	if err != nil {
		s.log.Errorf("Request SnapshotRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending SnapshotResponse %+v", response)
	return response, nil
}

func (s *Proxy) Restore(ctx context.Context, request *counter.RestoreRequest) (*counter.RestoreResponse, error) {
	s.log.Debugf("Received RestoreRequest %+v", request)
	header := request.Header
	partition := s.PartitionFor(header.PrimitiveID)

	conn, err := partition.Connect()
	if err != nil {
		return nil, err
	}

	client := counter.NewCounterServiceClient(conn)
	response, err := client.Restore(ctx, request)
	if err != nil {
		s.log.Errorf("Request RestoreRequest failed: %v", err)
		return nil, err
	}
	s.log.Debugf("Sending RestoreResponse %+v", response)
	return response, nil
}
