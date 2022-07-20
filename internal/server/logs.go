package server

func (h *Handler) LogsWriter(log string) {
	h.logs = append(h.logs, log+"\n")
}

func (h *Handler) LogsReader() (logs string) {
	for _, s := range h.logs {
		logs = logs + s
	}
	return logs
}
