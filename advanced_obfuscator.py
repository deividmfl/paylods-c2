#!/usr/bin/env python3
import os
import re
import glob
import random
import string
import hashlib

class AdvancedObfuscator:
    def __init__(self):
        self.string_map = {}
        self.class_map = {}
        self.method_map = {}
        self.variable_map = {}
        self.namespace_map = {}
        
    def generate_random_name(self, prefix="", length=8):
        """Gera nome aleatório baseado em hash"""
        chars = string.ascii_letters + string.digits
        base = ''.join(random.choices(chars, k=length))
        return f"{prefix}{base}" if prefix else base
    
    def generate_hash_name(self, original, type_prefix=""):
        """Gera nome baseado em hash do original"""
        hash_obj = hashlib.md5(original.encode())
        hash_hex = hash_obj.hexdigest()[:8]
        return f"{type_prefix}{hash_hex}"
    
    def obfuscate_strings(self, content):
        """Obfusca strings literais"""
        # Encontra strings entre aspas
        string_pattern = r'"([^"\\]*(\\.[^"\\]*)*)"'
        strings = re.findall(string_pattern, content)
        
        for string_match in strings:
            original_string = string_match[0]
            if len(original_string) > 3 and 'Apollo' in original_string:
                # Substitui Apollo por Phantom
                new_string = original_string.replace('Apollo', 'Phantom')
                content = content.replace(f'"{original_string}"', f'"{new_string}"')
                
        return content
    
    def obfuscate_identifiers(self, content):
        """Obfusca identificadores (classes, métodos, variáveis)"""
        
        # Padrões de identificadores C#
        patterns = {
            'namespace': r'namespace\s+([A-Za-z_][A-Za-z0-9_.]*)',
            'class': r'(public|private|internal|protected)?\s*(static|abstract|sealed)?\s*class\s+([A-Za-z_][A-Za-z0-9_]*)',
            'method': r'(public|private|internal|protected)\s+(static\s+)?(async\s+)?([A-Za-z_][A-Za-z0-9_]*\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*\(',
            'variable': r'(public|private|internal|protected)\s+(static\s+)?(readonly\s+)?([A-Za-z_][A-Za-z0-9_<>,\[\]]*)\s+([A-Za-z_][A-Za-z0-9_]*)\s*[=;]'
        }
        
        # Mapa de obfuscação específico para Apollo
        apollo_map = {
            'Apollo': 'X1a2b3c4',
            'ApolloInterop': 'Y2b3c4d5',
            'ApolloStructs': 'Z3c4d5e6',
            'TaskManager': 'A4d5e6f7',
            'FileManager': 'B5e6f7g8',
            'ProcessManager': 'C6f7g8h9',
            'InjectionManager': 'D7g8h9i0',
            'Config': 'E8h9i0j1',
            'PayloadUUID': 'F9i0j1k2',
            'JsonSerializer': 'G0j1k2l3',
            'MessageStore': 'H1k2l3m4',
            'Program': 'I2l3m4n5',
            'Main': 'J3m4n5o6',
            'Execute': 'K4n5o6p7',
            'HandleTasking': 'L5o6p7q8',
            'DeserializeToReceiverQueue': 'M6p7q8r9',
            'OnAsyncMessageReceived': 'N7q8r9s0',
            'Client_ConnectionEstablished': 'O8r9s0t1',
            '_jsonSerializer': 'P9s0t1u2',
            '_receiverQueue': 'Q0t1u2v3',
            '_connected': 'R1u2v3w4',
            '_cancellationToken': 'S2v3w4x5'
        }
        
        # Aplica mapeamento específico primeiro
        for original, obfuscated in apollo_map.items():
            content = re.sub(f'\\b{re.escape(original)}\\b', obfuscated, content)
        
        return content
    
    def add_junk_code(self, content):
        """Adiciona código lixo para confundir análise estática"""
        junk_methods = [
            '''
    private static void Xa1b2c3()
    {
        var x = DateTime.Now.Ticks;
        for(int i = 0; i < 10; i++)
        {
            x += i * 2;
        }
    }''',
            '''
    private static string Yb2c3d4()
    {
        return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy"));
    }''',
            '''
    private static void Zc3d4e5()
    {
        Thread.Sleep(Random.Next(1, 5));
        GC.Collect();
    }'''
        ]
        
        # Adiciona métodos lixo em classes
        class_pattern = r'(public\s+class\s+[A-Za-z0-9_]+\s*{)'
        matches = re.finditer(class_pattern, content)
        
        for match in matches:
            junk = random.choice(junk_methods)
            insertion_point = match.end()
            content = content[:insertion_point] + junk + content[insertion_point:]
            
        return content
    
    def scramble_control_flow(self, content):
        """Embaralha fluxo de controle básico"""
        # Adiciona declarações condicionais desnecessárias
        content = re.sub(
            r'(\s+)(return\s+[^;]+;)',
            r'\1if(DateTime.Now.Year > 2020) { \2 } else { return null; }',
            content
        )
        
        return content
    
    def obfuscate_file(self, file_path):
        """Obfusca um arquivo C# completamente"""
        print(f"Obfuscando avançado: {file_path}")
        
        try:
            with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
            
            # Aplica todas as técnicas de obfuscação
            content = self.obfuscate_strings(content)
            content = self.obfuscate_identifiers(content)
            content = self.add_junk_code(content)
            content = self.scramble_control_flow(content)
            
            # Remove comentários que possam revelar origem
            content = re.sub(r'//.*$', '', content, flags=re.MULTILINE)
            content = re.sub(r'/\*.*?\*/', '', content, flags=re.DOTALL)
            
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
                
        except Exception as e:
            print(f"Erro na obfuscação avançada de {file_path}: {e}")
    
    def obfuscate_directory(self, directory):
        """Obfusca todos os arquivos C# em um diretório"""
        cs_files = glob.glob(os.path.join(directory, '**', '*.cs'), recursive=True)
        
        for file_path in cs_files:
            if not any(skip in file_path for skip in ['Properties', 'obj', 'bin']):
                self.obfuscate_file(file_path)

if __name__ == "__main__":
    obfuscator = AdvancedObfuscator()
    agent_dir = "Payload_Types/phantom_apollo/agent_code"
    print("Iniciando obfuscação avançada...")
    obfuscator.obfuscate_directory(agent_dir)
    print("Obfuscação avançada concluída!")