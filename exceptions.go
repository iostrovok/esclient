package esclient

import (
	"strings"
)

type Code int

const (
	// https://godoc.org/google.golang.org/grpc/codes
	OK       Code = 0
	Unknown  Code = 2
	NotFound Code = 5
	Internal Code = 13
)

var exceptionsTypes map[string]Code

func (e *ErrorHandler) checkExceptionType() {

	e.code = OK

	if res, find := exceptionsTypes[strings.ToLower(e.esType)]; find {
		e.code = res
	} else if e.esStatus == 0 && e.httpStatusCode == 404 {
		e.code = NotFound
	} else if e.esStatus != 0 || e.httpStatusCode/200 != 1 {
		e.code = Unknown
	}

}

func init() {
	exceptionsTypes = map[string]Code{
		//"": OK,

		"action_not_found_transport_exception":         Internal,
		"action_transport_exception":                   Internal,
		"aggregation_execution_exception":              Internal,
		"aggregation_initialization_exception":         Internal,
		"alias_filter_parsing_exception":               Internal,
		"aliases_not_found_exception":                  Internal,
		"bind_http_exception":                          Internal,
		"bind_transport_exception":                     Internal,
		"blob_store_exception":                         Internal,
		"broadcast_shard_operation_failed_exception":   Internal,
		"circuit_breaking_exception":                   Internal,
		"cluster_block_exception":                      Internal,
		"concurrent_snapshot_execution_exception":      Internal,
		"connect_transport_exception":                  Internal,
		"coordination_state_rejected_exception":        Internal,
		"delay_recovery_exception":                     Internal,
		"dfs_phase_execution_exception":                Internal,
		"document_missing_exception":                   Internal,
		"document_source_missing_exception":            Internal,
		"elasticsearch_exception":                      Internal,
		"elasticsearch_generation_exception":           Internal,
		"elasticsearch_parse_exception":                Internal,
		"elasticsearch_security_exception":             Internal,
		"elasticsearch_timeout_exception":              Internal,
		"engine_creation_failure_exception":            Internal,
		"engine_exception":                             Internal,
		"execution_cancelled_exception":                Internal,
		"failed_node_exception":                        Internal,
		"failed_to_commit_cluster_state_exception":     Internal,
		"fetch_phase_execution_exception":              Internal,
		"flush_failed_engine_exception":                Internal,
		"gateway_exception":                            Internal,
		"general_script_exception":                     Internal,
		"http_exception":                               Internal,
		"http_on_transport_exception":                  Internal,
		"illegal_index_shard_state_exception":          Internal,
		"illegal_shard_routing_state_exception":        Internal,
		"incompatible_cluster_state_version_exception": Internal,
		"index_closed_exception":                       Internal,
		"index_creation_exception":                     Internal,
		"index_not_found_exception":                    Internal,
		"index_primary_shard_not_allocated_exception":  Internal,
		"index_shard_closed_exception":                 Internal,
		"index_shard_not_recovering_exception":         Internal,
		"index_shard_not_started_exception":            Internal,
		"index_shard_recovering_exception":             Internal,
		"index_shard_recovery_exception":               Internal,
		"index_shard_relocated_exception":              Internal,
		"index_shard_restore_exception":                Internal,
		"index_shard_restore_failed_exception":         Internal,
		"index_shard_snapshot_exception":               Internal,
		"index_shard_snapshot_failed_exception":        Internal,
		"index_shard_started_exception":                Internal,
		"index_template_missing_exception":             Internal,
		"invalid_aggregation_path_exception":           Internal,
		"invalid_alias_name_exception":                 Internal,
		"invalid_index_name_exception":                 Internal,
		"invalid_index_template_exception":             Internal,
		"invalid_snapshot_name_exception":              Internal,
		"invalid_type_name_exception":                  Internal,
		"mapper_exception":                             Internal,
		"mapper_parsing_exception":                     Internal,
		"master_not_discovered_exception":              Internal,
		"no_class_settings_exception":                  Internal,
		"no_longer_primary_shard_exception":            Internal,
		"no_node_available_exception":                  Internal,
		"no_shard_available_action_exception":          Internal,
		"no_such_node_exception":                       Internal,
		"no_such_remote_cluster_exception":             Internal,
		"node_closed_exception":                        Internal,
		"node_disconnected_exception":                  Internal,
		"node_not_connected_exception":                 Internal,
		"node_should_not_connect_exception":            Internal,
		"not_master_exception":                         Internal,
		"not_serializable_exception_wrapper":           Internal,
		"not_serializable_transport_exception":         Internal,
		"parsing_exception":                            Internal,
		"primary_missing_action_exception":             Internal,
		"process_cluster_event_timeout_exception":      Internal,
		"query_phase_execution_exception":              Internal,
		"query_shard_exception":                        Internal,
		"receive_timeout_transport_exception":          Internal,
		"recover_files_recovery_exception":             Internal,
		"recovery_engine_exception":                    Internal,
		"recovery_failed_exception":                    Internal,
		"reduce_search_phase_exception":                Internal,
		"refresh_failed_engine_exception":              Internal,
		"remote_transport_exception":                   Internal,
		"repository_exception":                         Internal,
		"repository_missing_exception":                 Internal,
		"repository_verification_exception":            Internal,
		"resource_already_exists_exception":            Internal,
		"resource_not_found_exception":                 Internal,
		"response_handler_failure_transport_exception": Internal,
		"retention_lease_already_exists_exception":     Internal,
		"retention_lease_not_found_exception":          Internal,
		"retry_on_primary_exception":                   Internal,
		"retry_on_replica_exception":                   Internal,
		"routing_exception":                            Internal,
		"routing_missing_exception":                    Internal,
		"script_exception":                             Internal,
		"search_context_exception":                     Internal,
		"search_context_missing_exception":             Internal,
		"search_exception":                             Internal,
		"search_parse_exception":                       Internal,
		"search_phase_execution_exception":             Internal,
		"search_source_builder_exception":              Internal,
		"send_request_transport_exception":             Internal,
		"settings_exception":                           Internal,
		"shard_lock_obtain_failed_exception":           Internal,
		"shard_not_found_exception":                    Internal,
		"shard_not_in_primary_mode_exception":          Internal,
		"snapshot_creation_exception":                  Internal,
		"snapshot_exception":                           Internal,
		"snapshot_failed_engine_exception":             Internal,
		"snapshot_in_progress_exception":               Internal,
		"snapshot_missing_exception":                   Internal,
		"snapshot_restore_exception":                   Internal,
		"status_exception":                             Internal,
		"strict_dynamic_mapping_exception":             Internal,
		"task_cancelled_exception":                     Internal,
		"timestamp_parsing_exception":                  Internal,
		"too_many_buckets_exception":                   Internal,
		"translog_corrupted_exception":                 Internal,
		"translog_exception":                           Internal,
		"transport_exception":                          Internal,
		"transport_serialization_exception":            Internal,
		"truncated_translog_exception":                 Internal,
		"type_missing_exception":                       Internal,
		"unavailable_shards_exception":                 Internal,
		"uncategorized_execution_exception":            Internal,
		"unknown_named_object_exception":               Internal,
		"version_conflict_engine_exception":            Internal,
	}
}
