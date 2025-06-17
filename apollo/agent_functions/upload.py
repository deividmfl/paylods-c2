from mythic_container import *

class UploadArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(
                name="file",
                type=ParameterType.File,
                description="File to upload",
                required=True,
            ),
            CommandParameter(
                name="remote_path",
                type=ParameterType.String,
                description="Remote path to save the file",
                required=True,
            )
        ]

    async def parse_arguments(self):
        if len(self.command_line) > 0:
            self.load_args_from_json_string(self.command_line)

class UploadCommand(CommandBase):
    cmd = "upload"
    needs_admin = False
    help_cmd = "upload"
    description = "Upload a file to the target"
    version = 1
    author = "@phantom"
    argument_class = UploadArguments
    attackmapping = ["T1105"]

    async def create_go_tasking(self, taskData: PTTaskMessageAllData) -> PTTaskCreateTaskingMessageResponse:
        response = PTTaskCreateTaskingMessageResponse(
            task_id=taskData.Task.ID,
            success=True,
        )
        return response

    async def process_response(self, task: PTTaskMessageAllData, response: any) -> PTTaskProcessResponseMessageResponse:
        resp = PTTaskProcessResponseMessageResponse(task_id=task.Task.ID, success=True)
        return resp