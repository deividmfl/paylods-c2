from mythic_container import MythicCommandBase, MythicTask
from mythic_container import TaskArguments


class PersistRemoveArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass


class PersistRemoveCommand(MythicCommandBase):
    cmd = "persist_remove"
    needs_admin = False
    help_cmd = "persist_remove"
    description = "Remove all persistence mechanisms"
    version = 1
    author = "@phantom"
    argument_class = PersistRemoveArguments
    attackmapping = ["T1070.004"]

    async def create_go_tasking(self, taskData: MythicTask) -> str:
        return "PersistenceHandler.RemovePersistence()"

    async def process_response(self, task: MythicTask, response: any) -> str:
        return f"Persistence removal: {response}"