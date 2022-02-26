package rocketmq

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"errors"
	"fmt"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/config"
)

var (
	ErrRocketmqPublishMsg         = errors.New("rocketmq publish msg error")
	ErrRocketmqValidPublishMsgTyp = errors.New("rocketmq publish msg error, invalid msg type")
)

const (
	metadataRocketmqTag           = "rocketmq-tag"
	metadataRocketmqKey           = "rocketmq-key"
	metadataRocketmqShardingKey   = "rocketmq-shardingkey"
	metadataRocketmqConsumerGroup = "rocketmq-consumerGroup"
	metadataRocketmqType          = "rocketmq-sub-type"
	metadataRocketmqExpression    = "rocketmq-sub-expression"
	metadataRocketmqBrokerName    = "rocketmq-broker-name"
)

type rocketMQMetaData struct {
	AccessProto string `mapstructure:"accessProto"`
	// rocketmq Credentials
	AccessKey  string `mapstructure:"accessKey"`
	SecretKey  string `mapstructure:"secretKey"`
	NameServer string `mapstructure:"nameServer"`
	GroupName  string `mapstructure:"groupName"`
	NameSpace  string `mapstructure:"nameSpace"`
	// consumer group rocketmq's subscribers
	ConsumerGroup     string `mapstructure:"consumerGroup"`
	ConsumerBatchSize int    `mapstructure:"consumerBatchSize"`
	// rocketmq's name server domain
	NameServerDomain string `mapstructure:"nameServerDomain"`
	// msg's content-type
	ContentType string `mapstructure:"content-type"`
	// retry times to connect rocketmq's broker
	Retries     int `mapstructure:"retries"`
	SendTimeOut int `mapstructure:"sendTimeOut"`
}

func getDefaultRocketMQMetaData() *rocketMQMetaData {
	return &rocketMQMetaData{
		AccessProto:       "",
		AccessKey:         "",
		SecretKey:         "",
		NameServer:        "",
		GroupName:         "",
		NameSpace:         "",
		ConsumerGroup:     "",
		ConsumerBatchSize: 0,
		NameServerDomain:  "",
		ContentType:       pubsub.DefaultCloudEventDataContentType,
		Retries:           3,
		SendTimeOut:       10,
	}
}

func (s *rocketMQMetaData) Decode(in interface{}) error {
	if err := config.Decode(in, &s); err != nil {
		return fmt.Errorf("decode failed. %w", err)
	}
	return nil
}

func parseRocketMQMetaData(metadata pubsub.Metadata) (*rocketMQMetaData, error) {
	rMetaData := getDefaultRocketMQMetaData()
	if metadata.Properties != nil {
		err := rMetaData.Decode(metadata.Properties)
		if err != nil {
			return nil, fmt.Errorf("rocketmq configuration error: %w", err)
		}
	}
	return rMetaData, nil
}
