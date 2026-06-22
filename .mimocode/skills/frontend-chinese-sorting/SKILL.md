---
name: frontend-chinese-sorting
description: Implement Chinese character sorting in frontend
---

# Frontend Chinese Character Sorting

This skill covers implementing proper sorting for mixed Chinese and ASCII characters in React applications.

## When to Use

- Displaying lists with both Chinese and English characters
- Need ASCII characters to appear before Chinese characters
- Want Chinese characters sorted by pinyin (phonetic)

## Problem

Default JavaScript sorting uses Unicode code points, which puts Chinese characters before ASCII characters. This is often not the desired behavior.

## Solution

Use `localeCompare` with Chinese locale support:

```javascript
const sortedItems = items.sort((a, b) => {
    const aIsAscii = a.charCodeAt(0) < 128
    const bIsAscii = b.charCodeAt(0) < 128
    
    // ASCII characters first
    if (aIsAscii && !bIsAscii) return -1
    if (!aIsAscii && bIsAscii) return 1
    
    // Same type: use locale compare
    return a.localeCompare(b, 'zh')
})
```

## Implementation in React

### With useMemo (Recommended)

```jsx
const sortedConfigs = useMemo(() => {
    const keys = Object.keys(data)
    const isAscii = (s) => s.charCodeAt(0) < 128
    
    return keys.sort((a, b) => {
        const aIsAscii = isAscii(a)
        const bIsAscii = isAscii(b)
        
        if (aIsAscii && !bIsAscii) return -1
        if (!aIsAscii && bIsAscii) return 1
        
        return a.localeCompare(b, 'zh')
    })
}, [data])
```

### With Regular Sort

```javascript
const sortChinese = (arr) => {
    const isAscii = (s) => s.charCodeAt(0) < 128
    
    return arr.sort((a, b) => {
        const aIsAscii = isAscii(a)
        const bIsAscii = isAscii(b)
        
        if (aIsAscii && !bIsAscii) return -1
        if (!aIsAscii && bIsAscii) return 1
        
        return a.localeCompare(b, 'zh')
    })
}
```

## Result Order

Example input: `["手游_PKM", "PC_AUG_X5", "手游_M7", "PC_K437", "PC_M7"]`

Sorted output: `["PC_AUG_X5", "PC_K437", "PC_M7", "手游_M7", "手游_PKM"]`

- ASCII characters: `PC_*` sorted alphabetically
- Chinese characters: `手游_*` sorted by pinyin

## Important Notes

- **Performance**: Use `useMemo` in React to avoid re-sorting on every render
- **Locale support**: `localeCompare('zh')` requires Chinese locale support in the browser
- **Edge cases**: Handle empty strings, numbers, and special characters
- **Accessibility**: Consider screen reader behavior for sorted lists

## Browser Support

- All modern browsers support `localeCompare` with locale parameter
- Chinese pinyin sorting works in Chrome, Firefox, Safari, Edge
- No polyfill needed for modern browser targets

## Related Files

- `server/src/App.jsx` - React component with sorted list
- `server/src/api.js` - API calls returning unsorted data
- `controller_MacroInterceptor.go` - Backend providing configuration data

## Testing

1. Test with mixed Chinese/ASCII list
2. Verify ASCII characters appear first
3. Check Chinese characters sorted by pinyin
4. Test with edge cases (empty strings, numbers)
5. Verify performance with large lists