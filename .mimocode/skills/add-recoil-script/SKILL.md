---
name: add-recoil-script
description: Add new weapon recoil compensation script
---

# Add Recoil Compensation Script

This skill covers adding new weapon recoil compensation scripts to the input2com system.

## When to Use

- Adding support for a new weapon's recoil pattern
- Updating existing recoil compensation data
- Testing new recoil patterns

## File Location

- **Directory**: `config/recoil/`
- **Format**: One `.txt` file per weapon
- **Naming**: Use weapon name (e.g., `AK47.txt`, `M4A1.txt`)

## File Format

Standard format (space-separated, 10ms slices):
```
dx dy
dx dy
dx dy
...
```

Where:
- `dx` = horizontal mouse movement (positive = right, negative = left)
- `dy` = vertical mouse movement (positive = down, negative = up)
- Each line represents 10ms of movement

## Custom Input Format

If using the custom recoil input feature:
```
period dx dy
```
Where:
- `period` = number of 10ms slices to repeat this movement
- Example: `3 0 5` means move down 5 pixels for 30ms (3 slices)

## Workflow

### 1. Create Script File

```bash
# Create new weapon script
echo "0 5
-1 3
2 4" > config/recoil/AK47.txt
```

### 2. Verify Auto-Loading

The system automatically loads all `.txt` files from `config/recoil/` at startup:
- Creates macro entry with weapon name
- Creates preset configuration
- Available in frontend UI

### 3. Test

1. Restart the program
2. Open web UI (port 9264)
3. Select the new weapon configuration
4. Test in application

## Important Notes

- **Auto-scanning**: Program scans `config/recoil/` directory at startup
- **No code changes needed**: Just add the file, system handles the rest
- **File format matters**: Must be space-separated `dx dy` pairs
- **10ms resolution**: Each line represents 10ms of movement
- **Chinese filenames supported**: System handles Unicode filenames

## Related Files

- `config/recoil/` - Recoil script directory
- `controller_MacroInterceptor.go` - Macro loading logic
- `configInit()` - Configuration initialization
- `NewMouseKeyboard_MacroInterceptor()` - Macro registration

## Example Scripts

See existing scripts in `config/recoil/` for reference:
- Simple weapons: `Pistol.txt` (minimal recoil)
- Complex weapons: `AK47.txt` (significant recoil pattern)
- Burst weapons: `M4A1.txt` (burst fire pattern)