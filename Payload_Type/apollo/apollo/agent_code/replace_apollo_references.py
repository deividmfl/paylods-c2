#!/usr/bin/env python3
"""
Script to systematically replace all Apollo references with Phantom throughout the codebase
"""
import os
import re
import glob

def replace_in_file(file_path, replacements):
    """Replace text in a file using the provided replacement mappings"""
    try:
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
        
        original_content = content
        
        # Apply all replacements
        for old_text, new_text in replacements.items():
            content = content.replace(old_text, new_text)
        
        # Only write if changes were made
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"Updated: {file_path}")
            return True
        return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    # Define all Apollo to Phantom replacements
    replacements = {
        # Class and namespace references
        'Apollo': 'Phantom',
        'apollo': 'phantom',
        'APOLLO': 'PHANTOM',
        
        # Specific namespace fixes
        'namespace SystemInterop.Enums.PhantomEnums': 'namespace SystemInterop.Enums.PhantomEnums',
        'namespace SystemInterop.Structs.PhantomStructs': 'namespace SystemInterop.Structs.PhantomStructs',
        
        # Assembly and product names
        'Product: Apollo': 'Product: Phantom',
        'AssemblyProduct("Apollo")': 'AssemblyProduct("Phantom")',
        'AssemblyTitle("Apollo")': 'AssemblyTitle("Phantom")',
        
        # Comments and documentation
        '// Apollo': '// Phantom',
        '* Apollo': '* Phantom',
        'Apollo agent': 'Phantom agent',
        'Apollo C2': 'Phantom C2',
        'Apollo framework': 'Phantom framework',
        
        # File references
        'Apollo.exe': 'Phantom.exe',
        'Apollo.dll': 'Phantom.dll',
        'Apollo.sln': 'Phantom.sln',
        'Apollo.csproj': 'Phantom.csproj',
        
        # Specific code patterns
        'ApolloInterop': 'PhantomInterop',
        'ApolloEnum': 'PhantomEnum',
        'ApolloStruct': 'PhantomStruct',
        'Apollo_': 'Phantom_',
        
        # Keep some system references intact (avoid breaking system calls)
        'Phantom.Agent.System': 'Agent.Phantom',
        'Agent.System': 'Agent.Phantom',
    }
    
    # File patterns to process
    file_patterns = [
        '**/*.cs',
        '**/*.csproj',
        '**/*.sln',
        '**/*.config',
        '**/*.xml'
    ]
    
    processed_files = 0
    updated_files = 0
    
    # Process all matching files
    for pattern in file_patterns:
        for file_path in glob.glob(pattern, recursive=True):
            # Skip binary files and specific directories
            if any(skip in file_path for skip in ['.git', 'bin', 'obj', '.vs']):
                continue
                
            processed_files += 1
            if replace_in_file(file_path, replacements):
                updated_files += 1
    
    print(f"\nProcessing complete:")
    print(f"Files processed: {processed_files}")
    print(f"Files updated: {updated_files}")

if __name__ == "__main__":
    main()