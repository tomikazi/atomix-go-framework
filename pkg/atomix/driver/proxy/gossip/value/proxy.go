// Code generated by atomix-go-sdk. DO NOT EDIT.
package value

import (
	"context"
	value "github.com/atomix/atomix-api/go/atomix/primitive/value"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/driver/proxy/gossip"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/errors"
	"github.com/atomix/atomix-go-sdk/pkg/atomix/logging"
	io "io"
)

// NewProxyServer creates a new ProxyServer
func NewProxyServer(server *gossip.Server) value.ValueServiceServer {
	return &ProxyServer{
		Server: server,
		log:    logging.GetLogger("atomix", "value"),
	}
}

type ProxyServer struct {
	*gossip.Server
	log logging.Logger
}

func (s *ProxyServer) Set(ctx context.Context, request *value.SetRequest) (*value.SetResponse, error) {
	s.log.Debugf("Received SetRequest %+v", request)
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := value.NewValueServiceClient(conn)
	ctx = partition.AddRequestHeaders(ctx, &request.Headers)
	response, err := client.Set(ctx, request)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition.AddResponseHeaders(&response.Headers)
	s.log.Debugf("Sending SetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Get(ctx context.Context, request *value.GetRequest) (*value.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := value.NewValueServiceClient(conn)
	ctx = partition.AddRequestHeaders(ctx, &request.Headers)
	response, err := client.Get(ctx, request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition.AddResponseHeaders(&response.Headers)
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Events(request *value.EventsRequest, srv value.ValueService_EventsServer) error {
	s.log.Debugf("Received EventsRequest %+v", request)
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	conn, err := partition.Connect()
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	client := value.NewValueServiceClient(conn)
	ctx := partition.AddRequestHeaders(srv.Context(), &request.Headers)
	stream, err := client.Events(ctx, request)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			s.log.Debugf("Finished EventsRequest %+v", request)
			return nil
		} else if err != nil {
			s.log.Errorf("Request EventsRequest failed: %v", err)
			return errors.Proto(err)
		}
		partition.AddResponseHeaders(&response.Headers)
		s.log.Debugf("Sending EventsResponse %+v", response)
		if err := srv.Send(response); err != nil {
			s.log.Errorf("Response EventsResponse failed: %v", err)
			return err
		}
	}
}
