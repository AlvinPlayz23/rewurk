Edit a file by exact find-and-replace. For multiple changes to one file, use edits[]; every edit is matched against the original file and the whole operation is atomic. The legacy top-level old_string/new_string form can also create missing files or delete content. For renames/moves use bash.

Matching is literal and exact: old_string must match the file byte-for-byte, including indentation (tabs vs spaces), blank lines, and trailing whitespace. Copy the text from read output (strip the line-number prefix) rather than retyping it. If old_string could appear more than once, include enough surrounding lines to make it unique, or set replace_all to replace every occurrence.

If the tool reports old_string not found, re-read the file at that location and copy more context exactly — never retry with a guessed or approximate match.
