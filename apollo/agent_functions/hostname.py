from mythic_container import *

class HostnameArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass

class HostnameCommand(CommandBase):
    cmd = "hostname"
    needs_admin = False
    help_cmd = "hostname"
    description = "Get the hostname of the target machine"
    version = 1
    author = "@phantom"
    argument_class = HostnameArguments
    attackmapping = ["T1082"]

    async def create_go_tasking(self, taskData: PTTaskMessageAllData) -> PTTaskCreateTaskingMessageResponse:
        response = PTTaskCreateTaskingMessageResponse(
            task_id=taskData.Task.ID,
            success=True,
        )
        return response

    async def process_response(self, task: PTTaskMessageAllData, response: any) -> PTTaskProcessResponseMessageResponse:
        resp = PTTaskProcessResponseMessageResponse(task_id=task.Task.ID, success=True)
        return resp