#!/usr/bin/env python3
"""
Assembly Rewriter for Phantom Apollo
Completely removes Apollo branding and references from .NET assemblies
"""

import os
import re
import sys
import hashlib
import shutil
from pathlib import Path

class AssemblyRewriter:
    def __init__(self):
        self.apollo_patterns = [
            # Assembly names and namespaces
            r'Apollo\.exe',
            r'Apollo,',
            r'Apollo\b',
            r'ApolloInterop',
            r'Apollo\.Agent',
            r'Apollo\.Classes',
            r'Apollo\.Structs',
            r'Apollo\.Enums',
            
            # Costura embedded resources
            r'costura\.apollointerop',
            r'costura\.apollo',
            
            # Metadata strings
            r'Apollo Agent',
            r'Apollo C2',
            r'Mythic Apollo',
            r'Apollo Framework',
            
            # HTTP headers and user agents
            r'Apollo/',
            r'apolloclient',
            
            # Class and method names
            r'ApolloTaskManager',
            r'ApolloAgent',
            r'ApolloStructs',
            r'ApolloEnums',
            r'ApolloLogonInformation'
        ]
        
        self.replacement_map = self.generate_replacements()
    
    def generate_replacements(self):
        """Generate consistent replacements for Apollo patterns"""
        replacements = {}
        
        # System-like replacements to blend in
        system_names = [
            'System', 'Service', 'Process', 'Network', 'Security',
            'Runtime', 'Core', 'Framework', 'Component', 'Module'
        ]
        
        for i, pattern in enumerate(self.apollo_patterns):
            # Use deterministic replacement based on pattern
            base_name = system_names[i % len(system_names)]
            
            if 'Apollo.exe' in pattern:
                replacements[pattern] = r'System.exe'
            elif 'ApolloInterop' in pattern:
                replacements[pattern] = r'SystemInterop'
            elif 'Apollo,' in pattern:
                replacements[pattern] = r'System,'
            elif 'Apollo\\b' in pattern:
                replacements[pattern] = r'System'
            elif 'costura.apollo' in pattern:
                replacements[pattern] = r'costura.system'
            elif 'Apollo Agent' in pattern:
                replacements[pattern] = r'System Service'
            elif 'Apollo C2' in pattern:
                replacements[pattern] = r'Network Service'
            elif 'Mythic Apollo' in pattern:
                replacements[pattern] = r'System Component'
            else:
                # Generate hash-based replacement
                hash_suffix = hashlib.md5(pattern.encode()).hexdigest()[:4]
                replacements[pattern] = f'{base_name}{hash_suffix}'
        
        return replacements
    
    def rewrite_csproj_files(self, directory):
        """Rewrite .csproj files to remove Apollo references"""
        print(f"[+] Rewriting .csproj files in {directory}")
        
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file.endswith('.csproj'):
                    file_path = os.path.join(root, file)
                    self.rewrite_xml_file(file_path)
    
    def rewrite_xml_file(self, file_path):
        """Rewrite XML configuration files"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # Replace Apollo references in XML
            xml_replacements = {
                '<AssemblyName>Apollo</AssemblyName>': '<AssemblyName>System</AssemblyName>',
                '<RootNamespace>Apollo</RootNamespace>': '<RootNamespace>System</RootNamespace>',
                '<Product>Apollo</Product>': '<Product>System Service</Product>',
                '<AssemblyTitle>Apollo</AssemblyTitle>': '<AssemblyTitle>System Service</AssemblyTitle>',
                '<AssemblyDescription>Apollo</AssemblyDescription>': '<AssemblyDescription>System Component</AssemblyDescription>',
                'Apollo.': 'System.',
                'apollo.': 'system.',
                'APOLLO': 'SYSTEM'
            }
            
            for old, new in xml_replacements.items():
                content = content.replace(old, new)
            
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"[+] Rewritten: {file_path}")
                
        except Exception as e:
            print(f"[-] Error rewriting {file_path}: {e}")
    
    def rewrite_cs_files(self, directory):
        """Rewrite C# source files to remove Apollo references"""
        print(f"[+] Rewriting C# files in {directory}")
        
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file.endswith('.cs'):
                    file_path = os.path.join(root, file)
                    self.rewrite_cs_file(file_path)
    
    def rewrite_cs_file(self, file_path):
        """Rewrite individual C# file"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # Replace namespace declarations
            content = re.sub(r'namespace\s+Apollo\b', 'namespace System', content)
            content = re.sub(r'namespace\s+ApolloInterop\b', 'namespace SystemInterop', content)
            
            # Replace class names
            content = re.sub(r'\bApollo\b', 'System', content)
            content = re.sub(r'\bApolloInterop\b', 'SystemInterop', content)
            content = re.sub(r'\bApolloAgent\b', 'SystemAgent', content)
            content = re.sub(r'\bApolloStructs\b', 'SystemStructs', content)
            content = re.sub(r'\bApolloEnums\b', 'SystemEnums', content)
            
            # Replace string literals
            content = re.sub(r'"Apollo[^"]*"', '"System"', content)
            content = re.sub(r"'Apollo[^']*'", "'System'", content)
            
            # Replace comments
            content = re.sub(r'//.*Apollo.*', '// System component', content)
            content = re.sub(r'/\*.*Apollo.*\*/', '/* System component */', content)
            
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"[+] Rewritten: {file_path}")
                
        except Exception as e:
            print(f"[-] Error rewriting {file_path}: {e}")
    
    def rewrite_assembly_info(self, directory):
        """Rewrite AssemblyInfo.cs files"""
        print(f"[+] Rewriting AssemblyInfo files in {directory}")
        
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file == 'AssemblyInfo.cs':
                    file_path = os.path.join(root, file)
                    self.rewrite_assembly_info_file(file_path)
    
    def rewrite_assembly_info_file(self, file_path):
        """Rewrite AssemblyInfo.cs file"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Replace assembly attributes
            replacements = {
                r'\[assembly: AssemblyTitle\("Apollo"\)\]': '[assembly: AssemblyTitle("System Service")]',
                r'\[assembly: AssemblyDescription\("Apollo[^"]*"\)\]': '[assembly: AssemblyDescription("System Component")]',
                r'\[assembly: AssemblyProduct\("Apollo"\)\]': '[assembly: AssemblyProduct("System Service")]',
                r'\[assembly: AssemblyCompany\("[^"]*"\)\]': '[assembly: AssemblyCompany("Microsoft Corporation")]',
                r'\[assembly: AssemblyCopyright\("Copyright.*Apollo.*"\)\]': '[assembly: AssemblyCopyright("Copyright © Microsoft Corporation")]'
            }
            
            for pattern, replacement in replacements.items():
                content = re.sub(pattern, replacement, content)
            
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"[+] Rewritten AssemblyInfo: {file_path}")
            
        except Exception as e:
            print(f"[-] Error rewriting AssemblyInfo {file_path}: {e}")
    
    def process_apollo_directory(self, apollo_dir):
        """Process entire Apollo directory structure"""
        print(f"[+] Processing Apollo directory: {apollo_dir}")
        
        if not os.path.exists(apollo_dir):
            print(f"[-] Directory not found: {apollo_dir}")
            return False
        
        # Rewrite source files
        self.rewrite_cs_files(apollo_dir)
        self.rewrite_csproj_files(apollo_dir)
        self.rewrite_assembly_info(apollo_dir)
        
        # Rename Apollo.sln if it exists
        sln_path = os.path.join(apollo_dir, 'Apollo.sln')
        if os.path.exists(sln_path):
            new_sln_path = os.path.join(apollo_dir, 'System.sln')
            shutil.move(sln_path, new_sln_path)
            print(f"[+] Renamed solution: Apollo.sln -> System.sln")
        
        print("[+] Apollo directory processing complete")
        return True

def main():
    """Main assembly rewriting function"""
    rewriter = AssemblyRewriter()
    
    # Process the Apollo agent code directory
    apollo_dirs = [
        "apollo/agent_code/Apollo",
        "apollo/agent_code",
        "Agent/Apollo",
        "Apollo"
    ]
    
    processed = False
    for apollo_dir in apollo_dirs:
        if os.path.exists(apollo_dir):
            rewriter.process_apollo_directory(apollo_dir)
            processed = True
            break
    
    if not processed:
        print("[-] No Apollo directory found for processing")
        return False
    
    print("[+] Assembly rewriting complete!")
    return True

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)