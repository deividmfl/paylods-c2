#!/usr/bin/env python3
"""
Advanced Code Obfuscator for Phantom Apollo
Applies comprehensive obfuscation to C# source files
"""
import os
import re
import random
import string
import hashlib
import base64
from pathlib import Path

class AdvancedObfuscator:
    def __init__(self):
        self.string_mapping = {}
        self.identifier_mapping = {}
        self.method_mapping = {}
        self.class_mapping = {}
        
    def generate_random_name(self, prefix="", length=8):
        """Gera nome aleatório baseado em hash"""
        chars = string.ascii_letters + string.digits
        random_part = ''.join(random.choices(chars, k=length))
        return f"{prefix}{random_part}"
    
    def generate_hash_name(self, original, type_prefix=""):
        """Gera nome baseado em hash do original"""
        hash_obj = hashlib.md5(original.encode())
        hash_hex = hash_obj.hexdigest()[:8]
        return f"{type_prefix}{hash_hex}"
    
    def obfuscate_strings(self, content):
        """Obfusca strings literais"""
        # Pattern para strings literais C#
        string_pattern = r'"([^"\\]|\\.)*"'
        
        def replace_string(match):
            original_string = match.group(0)
            string_content = original_string[1:-1]  # Remove quotes
            
            if len(string_content) < 3:  # Skip very short strings
                return original_string
            
            # Encode string
            encoded = base64.b64encode(string_content.encode()).decode()
            obfuscated_var = self.generate_random_name("str_", 6)
            
            # Store mapping for later injection
            self.string_mapping[obfuscated_var] = encoded
            
            # Return obfuscated call
            return f'System.Text.Encoding.UTF8.GetString(System.Convert.FromBase64String("{encoded}"))'
        
        return re.sub(string_pattern, replace_string, content)
    
    def obfuscate_identifiers(self, content):
        """Obfusca identificadores (classes, métodos, variáveis)"""
        # Classes
        class_pattern = r'\bclass\s+([A-Za-z_][A-Za-z0-9_]*)'
        def replace_class(match):
            original = match.group(1)
            if original not in self.class_mapping:
                self.class_mapping[original] = self.generate_hash_name(original, "C")
            return f"class {self.class_mapping[original]}"
        
        content = re.sub(class_pattern, replace_class, content)
        
        # Métodos públicos
        method_pattern = r'\b(public|private|protected|internal)\s+(\w+\s+)?(\w+)\s*\('
        def replace_method(match):
            visibility = match.group(1)
            return_type = match.group(2) or ""
            method_name = match.group(3)
            
            if method_name not in ['Main', 'ToString', 'Equals', 'GetHashCode']:  # Preserve special methods
                if method_name not in self.method_mapping:
                    self.method_mapping[method_name] = self.generate_hash_name(method_name, "M")
                method_name = self.method_mapping[method_name]
            
            return f"{visibility} {return_type}{method_name}("
        
        content = re.sub(method_pattern, replace_method, content)
        
        return content
    
    def add_junk_code(self, content):
        """Adiciona código lixo para confundir análise estática"""
        junk_methods = []
        
        for i in range(random.randint(3, 7)):
            method_name = self.generate_random_name("Junk", 8)
            junk_code = f"""
        private static void {method_name}()
        {{
            var {self.generate_random_name("var", 4)} = new Random().Next();
            if ({self.generate_random_name("var", 4)} > 999999999)
            {{
                Console.WriteLine("Never executed");
            }}
            Thread.Sleep(0);
        }}"""
            junk_methods.append(junk_code)
        
        # Insert junk methods before the last closing brace
        last_brace_pos = content.rfind('}')
        if last_brace_pos != -1:
            content = content[:last_brace_pos] + '\n'.join(junk_methods) + '\n' + content[last_brace_pos:]
        
        return content
    
    def scramble_control_flow(self, content):
        """Embaralha fluxo de controle básico"""
        # Add fake conditional blocks
        fake_conditions = [
            "if (DateTime.Now.Year > 2030) return;",
            "if (Environment.ProcessorCount < 0) Environment.Exit(0);",
            "if (new Random().Next() > int.MaxValue) throw new Exception();",
        ]
        
        # Insert at random positions in methods
        method_bodies = re.finditer(r'\{([^{}]*(?:\{[^{}]*\}[^{}]*)*)\}', content)
        
        for match in method_bodies:
            body = match.group(1)
            if len(body.strip()) > 50:  # Only add to substantial methods
                fake_condition = random.choice(fake_conditions)
                new_body = body.replace('{', '{\n    ' + fake_condition + '\n', 1)
                content = content.replace(match.group(0), '{' + new_body + '}')
        
        return content
    
    def obfuscate_file(self, file_path):
        """Obfusca um arquivo C# completamente"""
        print(f"[+] Obfuscating {file_path}")
        
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Apply obfuscation layers
        content = self.obfuscate_strings(content)
        content = self.obfuscate_identifiers(content)
        content = self.add_junk_code(content)
        content = self.scramble_control_flow(content)
        
        # Write obfuscated file
        obfuscated_path = file_path.replace('.cs', '_obfuscated.cs')
        with open(obfuscated_path, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"[+] Created obfuscated file: {obfuscated_path}")
        return obfuscated_path
    
    def obfuscate_directory(self, directory):
        """Obfusca todos os arquivos C# em um diretório"""
        cs_files = Path(directory).rglob("*.cs")
        obfuscated_files = []
        
        for cs_file in cs_files:
            # Skip certain files
            if any(skip in str(cs_file) for skip in ['AssemblyInfo', 'GlobalAssemblyInfo', 'TemporaryGeneratedFile']):
                continue
                
            obfuscated_file = self.obfuscate_file(str(cs_file))
            obfuscated_files.append(obfuscated_file)
        
        return obfuscated_files

def main():
    import sys
    
    if len(sys.argv) != 2:
        print("Usage: python advanced_obfuscator.py <directory>")
        sys.exit(1)
    
    directory = sys.argv[1]
    if not os.path.exists(directory):
        print(f"Error: Directory {directory} not found")
        sys.exit(1)
    
    obfuscator = AdvancedObfuscator()
    obfuscated_files = obfuscator.obfuscate_directory(directory)
    
    print(f"\n[+] Obfuscation complete!")
    print(f"[+] Files processed: {len(obfuscated_files)}")
    for file in obfuscated_files:
        print(f"    - {file}")

if __name__ == "__main__":
    main()