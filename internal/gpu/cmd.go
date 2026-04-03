package gpu

import "github.com/fumiama/gozel/ze"

// CommandListCreate creates a command list on the default device.
func CommandListCreate() (ze.CommandListHandle, error) {
	return ctx.CommandListCreate(dev)
}

// ExecuteCommandLists submits the command list for execution on the command queue.
func ExecCommandLists(hCommandList ...ze.CommandListHandle) error {
	return q.ExecuteCommandLists(hCommandList...)
}
