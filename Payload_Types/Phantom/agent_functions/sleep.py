from mythic_container.MythicCommandBase import *
from mythic_container.MythicRPC import *

class SleepArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(
                name="time",
                type=ParameterType.Number,
                description="Sleep time in seconds",
                parameter_group_info=[ParameterGroupInfo(
                    required=True,
                    group_name="Default"
                )]
            )
        ]

    async def parse_arguments(self):
        if len(self.command_line) == 0:
            raise ValueError("Must supply sleep time")
        try:
            sleep_time = int(self.command_line)
            self.add_arg("time", sleep_time)
        except ValueError:
            raise ValueError("Sleep time must be a number")

class SleepCommand(CommandBase):
    cmd = "sleep"
    needs_admin = False
    help_cmd = "sleep [seconds]"
    description = "Change the sleep interval for the agent"
    version = 1
    author = "@phantom"
    argument_class = SleepArguments
    attackmapping = []
    
    async def create_go_tasking(self, taskData: PTTaskMessageAllData) -> PTTaskCreateTaskingMessageResponse:
        response = PTTaskCreateTaskingMessageResponse(
            TaskID=taskData.Task.ID,
            Success=True,
        )
        return response

    async def process_response(self, task: PTTaskMessageAllData, response: any) -> PTTaskProcessResponseMessageResponse:
        resp = PTTaskProcessResponseMessageResponse(TaskID=task.Task.ID, Success=True)
        return resp