// Code generated by atomix-go-framework. DO NOT EDIT.

// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package election

import (
	"context"
	election "github.com/atomix/atomix-api/go/atomix/primitive/election"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	storage "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
)

const Type = "Election"

const (
	enterOp    storage.OperationID = 1
	withdrawOp storage.OperationID = 2
	anointOp   storage.OperationID = 3
	promoteOp  storage.OperationID = 4
	evictOp    storage.OperationID = 5
	getTermOp  storage.OperationID = 6
	eventsOp   storage.OperationID = 7
)

var log = logging.GetLogger("atomix", "proxy", "election")

// NewProxyServer creates a new ProxyServer
func NewProxyServer(client *rsm.Client, readSync bool) election.LeaderElectionServiceServer {
	return &ProxyServer{
		Client:   client,
		readSync: readSync,
	}
}

type ProxyServer struct {
	*rsm.Client
	readSync bool
	log      logging.Logger
}

func (s *ProxyServer) Enter(ctx context.Context, request *election.EnterRequest) (*election.EnterResponse, error) {
	log.Debugf("Received EnterRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request EnterRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request EnterRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoCommand(ctx, enterOp, input)
	if err != nil {
		log.Debugf("Request EnterRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.EnterResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request EnterRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending EnterResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) Withdraw(ctx context.Context, request *election.WithdrawRequest) (*election.WithdrawResponse, error) {
	log.Debugf("Received WithdrawRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request WithdrawRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request WithdrawRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoCommand(ctx, withdrawOp, input)
	if err != nil {
		log.Debugf("Request WithdrawRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.WithdrawResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request WithdrawRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending WithdrawResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) Anoint(ctx context.Context, request *election.AnointRequest) (*election.AnointResponse, error) {
	log.Debugf("Received AnointRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request AnointRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request AnointRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoCommand(ctx, anointOp, input)
	if err != nil {
		log.Debugf("Request AnointRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.AnointResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request AnointRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending AnointResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) Promote(ctx context.Context, request *election.PromoteRequest) (*election.PromoteResponse, error) {
	log.Debugf("Received PromoteRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request PromoteRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request PromoteRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoCommand(ctx, promoteOp, input)
	if err != nil {
		log.Debugf("Request PromoteRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.PromoteResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request PromoteRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending PromoteResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) Evict(ctx context.Context, request *election.EvictRequest) (*election.EvictResponse, error) {
	log.Debugf("Received EvictRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request EvictRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request EvictRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoCommand(ctx, evictOp, input)
	if err != nil {
		log.Debugf("Request EvictRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.EvictResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request EvictRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending EvictResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) GetTerm(ctx context.Context, request *election.GetTermRequest) (*election.GetTermResponse, error) {
	log.Debugf("Received GetTermRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request GetTermRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(ctx, serviceInfo)
	if err != nil {
		log.Errorf("Request GetTermRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := service.DoQuery(ctx, getTermOp, input, s.readSync)
	if err != nil {
		log.Debugf("Request GetTermRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &election.GetTermResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request GetTermRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending GetTermResponse %.250s", response)
	return response, nil
}

func (s *ProxyServer) Events(request *election.EventsRequest, srv election.LeaderElectionService_EventsServer) error {
	log.Debugf("Received EventsRequest %.250s", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	serviceInfo := storage.ServiceInfo{
		Type:      storage.ServiceType(Type),
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	service, err := partition.GetService(srv.Context(), serviceInfo)
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	err = service.DoCommandStream(srv.Context(), eventsOp, input, stream)
	if err != nil {
		log.Debugf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	ch := make(chan streams.Result)
	go func() {
		defer close(ch)
		for {
			result, ok := stream.Receive()
			if !ok {
				return
			}
			ch <- result
		}
	}()

	for result := range ch {
		if result.Failed() {
			if result.Error == context.Canceled {
				break
			}
			log.Debugf("Request EventsRequest failed: %v", result.Error)
			return errors.Proto(result.Error)
		}

		response := &election.EventsResponse{}
		err = proto.Unmarshal(result.Value.([]byte), response)
		if err != nil {
			log.Errorf("Request EventsRequest failed: %v", err)
			return errors.Proto(err)
		}

		log.Debugf("Sending EventsResponse %.250s", response)
		if err = srv.Send(response); err != nil {
			log.Warnf("Response EventsResponse failed: %v", err)
			return err
		}
	}
	log.Debugf("Finished EventsRequest %.250s", request)
	return nil
}
