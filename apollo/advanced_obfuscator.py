#!/usr/bin/env python3
"""
Advanced Code Obfuscator for Phantom Apollo
Applies comprehensive obfuscation to C# source files
"""

import os
import re
import hashlib
import random
import string
from pathlib import Path

class AdvancedObfuscator:
    def __init__(self):
        self.identifier_map = {}
        self.string_map = {}
        self.method_map = {}
        
        # Pre-defined obfuscated names based on original Apollo mappings
        self.core_mappings = {
            # Namespaces
            'Apollo': 'X1a2b3c4',
            'ApolloInterop': 'PhantomInterop',
            
            # Core classes
            'TaskManager': 'CommandProcessor',
            'FileManager': 'DataHandler',
            'Agent': 'Runtime',
            'Config': 'Settings',
            
            # Common variables
            '_jsonSerializer': '_dataSerializer',
            '_receiverQueue': '_msgRecvQueue',
            '_senderQueue': '_msgSendQueue',
            'PayloadUUID': 'AgentIdentifier',
            'Main': 'J3m4n5o6',
            
            # Methods
            'DeserializeToReceiverQueue': 'HandleIncomingData',
            'HandleTasking': 'ProcessCommand',
            'SendResult': 'TransmitResponse'
        }

    def generate_random_name(self, prefix="", length=8):
        """Gera nome aleatório baseado em hash"""
        chars = string.ascii_letters + string.digits
        return prefix + ''.join(random.choice(chars) for _ in range(length))

    def generate_hash_name(self, original, type_prefix=""):
        """Gera nome baseado em hash do original"""
        if original in self.core_mappings:
            return self.core_mappings[original]
        
        hash_obj = hashlib.md5(original.encode())
        hash_hex = hash_obj.hexdigest()[:8]
        return f"{type_prefix}{hash_hex.upper()}"

    def obfuscate_strings(self, content):
        """Obfusca strings literais"""
        # Find string literals
        string_pattern = r'"([^"\\]|\\.)*"'
        strings = re.findall(string_pattern, content)
        
        for string_literal in strings:
            if len(string_literal) > 6 and not any(x in string_literal.lower() for x in ['http', 'www', '.exe', '.dll']):
                # Generate obfuscated version
                original = string_literal[1:-1]  # Remove quotes
                encoded = ''.join([f"\\x{ord(c):02x}" for c in original])
                obfuscated = f'"{encoded}"'
                content = content.replace(string_literal, obfuscated)
        
        return content

    def obfuscate_identifiers(self, content):
        """Obfusca identificadores (classes, métodos, variáveis)"""
        # Class declarations
        class_pattern = r'\bclass\s+(\w+)'
        for match in re.finditer(class_pattern, content):
            original = match.group(1)
            if original not in self.identifier_map:
                self.identifier_map[original] = self.generate_hash_name(original, "C")
            content = re.sub(rf'\b{original}\b', self.identifier_map[original], content)
        
        # Method declarations
        method_pattern = r'\b(public|private|protected|internal)\s+(?:static\s+)?(?:async\s+)?(?:\w+\s+)?(\w+)\s*\('
        for match in re.finditer(method_pattern, content):
            original = match.group(2)
            if original not in ['Main', 'ToString', 'GetHashCode', 'Equals'] and not original.startswith('get_') and not original.startswith('set_'):
                if original not in self.method_map:
                    self.method_map[original] = self.generate_hash_name(original, "M")
                content = re.sub(rf'\b{original}\b', self.method_map[original], content)
        
        # Apply core mappings
        for original, obfuscated in self.core_mappings.items():
            content = re.sub(rf'\b{original}\b', obfuscated, content)
        
        return content

    def add_junk_code(self, content):
        """Adiciona código lixo para confundir análise estática"""
        junk_methods = [
            '''
    private static void Z9x8c7v6()
    {
        var tmp = DateTime.Now.Ticks;
        for (int i = 0; i < tmp % 100; i++)
        {
            var dummy = i * 2;
        }
    }''',
            '''
    private static bool B3n4m5l6()
    {
        return Environment.TickCount % 2 == 0;
    }''',
            '''
    private static void Q1w2e3r4()
    {
        var random = new Random();
        var array = new int[random.Next(10, 100)];
        for (int i = 0; i < array.Length; i++)
        {
            array[i] = random.Next();
        }
    }'''
        ]
        
        # Insert junk code before class closing brace
        for junk in junk_methods:
            if random.random() > 0.7:  # 30% chance to add each junk method
                content = content.replace('    }', junk + '\n    }', 1)
        
        return content

    def scramble_control_flow(self, content):
        """Embaralha fluxo de controle básico"""
        # Add unnecessary try-catch blocks
        simple_statements = re.findall(r'(\s+)(\w+\.\w+\([^)]*\);)', content)
        for indent, statement in simple_statements:
            if random.random() > 0.8:  # 20% chance
                wrapped = f'''{indent}try
{indent}{{
{indent}    {statement}
{indent}}}
{indent}catch (Exception) {{ }}'''
                content = content.replace(indent + statement, wrapped)
        
        return content

    def obfuscate_file(self, file_path):
        """Obfusca um arquivo C# completamente"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            print(f"[+] Obfuscating {file_path}")
            
            # Apply obfuscation techniques
            content = self.obfuscate_identifiers(content)
            content = self.obfuscate_strings(content)
            content = self.add_junk_code(content)
            content = self.scramble_control_flow(content)
            
            # Remove comments
            content = re.sub(r'//.*$', '', content, flags=re.MULTILINE)
            content = re.sub(r'/\*.*?\*/', '', content, flags=re.DOTALL)
            
            # Backup original
            backup_path = file_path + '.original'
            if not os.path.exists(backup_path):
                with open(backup_path, 'w', encoding='utf-8') as f:
                    f.write(content)
            
            # Write obfuscated version
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            
            return True
        except Exception as e:
            print(f"[-] Error obfuscating {file_path}: {e}")
            return False

    def obfuscate_directory(self, directory):
        """Obfusca todos os arquivos C# em um diretório"""
        cs_files = []
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file.endswith('.cs') and not file.endswith('.original'):
                    cs_files.append(os.path.join(root, file))
        
        print(f"[+] Found {len(cs_files)} C# files to obfuscate")
        
        success_count = 0
        for file_path in cs_files:
            if self.obfuscate_file(file_path):
                success_count += 1
        
        print(f"[+] Successfully obfuscated {success_count}/{len(cs_files)} files")
        return success_count == len(cs_files)

def main():
    print("[+] Phantom Apollo Advanced Obfuscator")
    
    # Obfuscate agent code directory
    agent_code_dir = "agent_code"
    if os.path.exists(agent_code_dir):
        obfuscator = AdvancedObfuscator()
        obfuscator.obfuscate_directory(agent_code_dir)
    else:
        print(f"[-] Directory {agent_code_dir} not found")
        return False
    
    print("[+] Advanced obfuscation completed")
    return True

if __name__ == "__main__":
    main()