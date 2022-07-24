package server

func (h *Handler) LogsWriter(log string) {
	h.mut.Lock()
	h.logs = append(h.logs, log+"\n")
	h.mut.Unlock()
}

func (h *Handler) LogsReader() (logs string) {
	h.mut.Lock()
	for _, s := range h.logs {
		logs = logs + s
	}
	h.mut.Unlock()
	return logs
}
