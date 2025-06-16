from mythic_container.MythicCommandBase import *
from mythic_container.MythicRPC import *

class UploadArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(
                name="file",
                type=ParameterType.File,
                description="File to upload",
                parameter_group_info=[ParameterGroupInfo(
                    required=True,
                    group_name="Default"
                )]
            ),
            CommandParameter(
                name="path",
                type=ParameterType.String,
                description="Destination path on target",
                parameter_group_info=[ParameterGroupInfo(
                    required=True,
                    group_name="Default"
                )]
            )
        ]

    async def parse_arguments(self):
        if len(self.command_line) == 0:
            raise ValueError("Must supply destination path")
        self.add_arg("path", self.command_line)

class UploadCommand(CommandBase):
    cmd = "upload"
    needs_admin = False
    help_cmd = "upload [destination_path]"
    description = "Upload a file to the target"
    version = 1
    author = "@phantom"
    argument_class = UploadArguments
    attackmapping = ["T1105"]
    
    async def create_go_tasking(self, taskData: PTTaskMessageAllData) -> PTTaskCreateTaskingMessageResponse:
        response = PTTaskCreateTaskingMessageResponse(
            TaskID=taskData.Task.ID,
            Success=True,
        )
        return response

    async def process_response(self, task: PTTaskMessageAllData, response: any) -> PTTaskProcessResponseMessageResponse:
        resp = PTTaskProcessResponseMessageResponse(TaskID=task.Task.ID, Success=True)
        return resp