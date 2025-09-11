package model

type KafkaConfig struct {
	BootstrapServersValue    string `yaml:"bootstrap_servers_value"`
	ResponseMessageTopicName string `yaml:"response_message_topic_name"`
	RequestMessageTopicName  string `yaml:"request_message_topic_name"`
	RequestMessageGroupId    string `yaml:"request_message_group_id"`
}

func (KafkaConfig) SectionName() string {
	return "kafka"
}
