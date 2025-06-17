#!/usr/bin/env python3
import os
import re
import glob

# Mapeamento de obfuscação para nomes de classes, namespaces e variáveis
obfuscation_map = {
    # Namespaces principais
    'Apollo': 'Phantom',
    'ApolloInterop': 'PhantomInterop',
    'ApolloStructs': 'PhantomStructs',
    'ApolloEnums': 'PhantomEnums',
    
    # Classes principais
    'JsonSerializer': 'JsonHandler',
    'Program': 'Runtime',
    'TaskManager': 'CommandProcessor',
    'FileManager': 'DataHandler',
    'IdentityManager': 'UserContext',
    'C2ProfileManager': 'CommHandler',
    'PeerManager': 'NodeHandler',
    'SocksManager': 'ProxyHandler',
    'RpfwdManager': 'TunnelHandler',
    'ProcessManager': 'ProcHandler',
    'InjectionManager': 'CodeInjector',
    'KerberosTicketManager': 'TicketHandler',
    
    # Variáveis comuns
    '_jsonSerializer': '_dataSerializer',
    '_receiverEvent': '_msgRecvEvent',
    '_receiverQueue': '_msgRecvQueue',
    'MessageStore': 'DataStore',
    '_connected': '_connectionActive',
    '_senderQueue': '_msgSendQueue',
    '_sendAction': '_transmitAction',
    '_cancellationToken': '_stopToken',
    '_senderEvent': '_msgSendEvent',
    '_complete': '_taskComplete',
    '_completed': '_isFinished',
    '_flushMessages': '_flushData',
    
    # Métodos comuns
    'DeserializeToReceiverQueue': 'HandleIncomingData',
    'OnAsyncMessageReceived': 'ProcessReceivedMessage',
    'OnAsyncMessageSent': 'ProcessSentMessage',
    'Client_ConnectionEstablished': 'OnConnectionReady',
    'Client_Disconnect': 'OnConnectionClosed',
    
    # Configurações
    'Config.PayloadUUID': 'Settings.AgentIdentifier',
    'Config.EgressProfiles': 'Settings.CommProfiles',
    'Config.StagingRSAPrivateKey': 'Settings.CryptoKey',
    
    # Outros
    'IMythicMessage': 'ICommandMessage',
    'IPCChunkedData': 'DataChunk',
    'ChunkedMessageStore': 'ChunkStore',
    'NamedPipeMessageArgs': 'PipeMessageData',
    'ChunkMessageEventArgs': 'ChunkEventData',
}

def obfuscate_file(file_path):
    """Obfusca um arquivo aplicando as substituições do mapeamento"""
    print(f"Obfuscando: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
        content = f.read()
    
    # Aplica todas as substituições do mapeamento
    for original, obfuscated in obfuscation_map.items():
        # Usa regex para substituições mais precisas
        # Substitui declarações de namespace
        content = re.sub(f'namespace {re.escape(original)}\\b', f'namespace {obfuscated}', content)
        # Substitui using statements
        content = re.sub(f'using {re.escape(original)}\\b', f'using {obfuscated}', content)
        # Substitui nomes de classes e variáveis
        content = re.sub(f'\\b{re.escape(original)}\\b', obfuscated, content)
    
    # Substitui comentários identificadores do Apollo
    content = re.sub(r'// Apollo', '// Phantom', content)
    content = re.sub(r'/\* Apollo', '/* Phantom', content)
    content = re.sub(r'Apollo Agent', 'Phantom Agent', content)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)

def obfuscate_directory(directory):
    """Obfusca todos os arquivos C# em um diretório recursivamente"""
    cs_files = glob.glob(os.path.join(directory, '**', '*.cs'), recursive=True)
    config_files = glob.glob(os.path.join(directory, '**', '*.config'), recursive=True)
    project_files = glob.glob(os.path.join(directory, '**', '*.csproj'), recursive=True)
    
    all_files = cs_files + config_files + project_files
    
    for file_path in all_files:
        try:
            obfuscate_file(file_path)
        except Exception as e:
            print(f"Erro ao obfuscar {file_path}: {e}")

if __name__ == "__main__":
    agent_code_dir = "Payload_Types/phantom_apollo/agent_code"
    print(f"Iniciando obfuscação do diretório: {agent_code_dir}")
    obfuscate_directory(agent_code_dir)
    print("Obfuscação concluída!")