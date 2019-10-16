package api

// Merge takes two PipelineDocuments and merge them
func (oldDoc *PipelineDocument) Merge(newDoc *PipelineDocument) {
	// Update some fields of oldDoc to match newDoc now
	if newDoc.AppConfig != nil {
		oldDoc.AppConfig = newDoc.AppConfig
	}
	if newDoc.Description != "" {
		oldDoc.Description = newDoc.Description
	}
	if newDoc.ExecutionEngine != "" {
		oldDoc.ExecutionEngine = newDoc.ExecutionEngine
	}
	if newDoc.Parallel != nil {
		oldDoc.Parallel = newDoc.Parallel
	}
	if newDoc.LimitConcurrent != nil {
		oldDoc.LimitConcurrent = newDoc.LimitConcurrent
	}
	if newDoc.KeepWaitingPipelines != nil {
		oldDoc.KeepWaitingPipelines = newDoc.KeepWaitingPipelines
	}
	if newDoc.Triggers != nil {
		oldDoc.Triggers = newDoc.Triggers
	}
	if newDoc.Parameters != nil {
		for _, newParam := range newDoc.Parameters {
			found := false
			for idx, oldParam := range oldDoc.Parameters {
				if oldParam.Name == newParam.Name {
					oldDoc.Parameters[idx] = newParam
					found = true
					continue
				}
			}
			if !found {
				oldDoc.Parameters = append(oldDoc.Parameters, newParam)
			}
		}
	}
	if newDoc.Stages != nil {
		for _, newStage := range newDoc.Stages {
			found := false
			for idx, oldStage := range oldDoc.Stages {
				if oldStage.RefID == newStage.RefID {
					oldDoc.Stages[idx] = newStage
					found = true
					continue
				}
			}
			if !found {
				oldDoc.Stages = append(oldDoc.Stages, newStage)
			}
		}
	}
}
