# Task Drag and Drop Guide

## Overview

The Discovery Tree frontend now supports drag-and-drop functionality for reordering tasks. You can drag tasks to change their position within the same parent or move them to different parents.

## How to Use Drag and Drop

### 1. Dragging Tasks

- **Click and hold** any task in the tree
- The task will become semi-transparent and slightly rotated to indicate it's being dragged
- **Drag the task** to your desired location

### 2. Drop Positions

When dragging a task over another task, you'll see visual indicators for three possible drop positions:

#### **Before** (Top 25% of target task)
- A **blue line appears above** the target task
- Drops the task **before** the target task (same parent)

#### **After** (Bottom 25% of target task)  
- A **blue line appears below** the target task
- Drops the task **after** the target task (same parent)

#### **Child** (Middle 50% of target task)
- The target task gets a **blue dashed border**
- Drops the task **as a child** of the target task

### 3. Visual Feedback

- **Dragging task**: Semi-transparent with slight rotation
- **Valid drop targets**: Highlighted with blue background
- **Drop indicators**: Blue lines (before/after) or dashed border (child)
- **Invalid drops**: No visual feedback (e.g., dropping on itself)

## Examples

### Reordering Siblings
```
Before:
📋 Project
  ├─ Task A
  ├─ Task B  ← Drag this
  └─ Task C

After (dropping Task B after Task C):
📋 Project
  ├─ Task A
  ├─ Task C
  └─ Task B
```

### Moving to Different Parent
```
Before:
📋 Project A
  └─ Task 1  ← Drag this
📋 Project B
  └─ Task 2

After (dropping Task 1 as child of Project B):
📋 Project A
📋 Project B
  ├─ Task 2
  └─ Task 1
```

### Creating Subtasks
```
Before:
📋 Project
  ├─ Task A
  └─ Task B  ← Drag this

After (dropping Task B as child of Task A):
📋 Project
  └─ Task A
      └─ Task B
```

## Limitations

- Cannot drag tasks while editing (inline editing mode)
- Cannot drop a task onto itself
- Cannot create circular dependencies (parent cannot become child of its descendant)
- Drag and drop is disabled during loading states

## Keyboard Accessibility

While drag-and-drop is primarily a mouse/touch interaction, you can still reorder tasks using:
- **Right-click context menu** → "Move" options (if implemented)
- **Cut/Copy/Paste** functionality (if implemented)

## Technical Details

- Uses HTML5 Drag and Drop API
- Validates moves using TreeService before API calls
- Optimistic updates with error rollback
- Automatic tree refresh after successful moves
- Position calculations based on existing task positions

## Troubleshooting

**Task won't drag:**
- Make sure you're not in edit mode (press Escape to exit)
- Check that the task isn't in a loading state

**Drop doesn't work:**
- Ensure you're dropping on a valid target
- Check browser console for any error messages
- Verify the backend API is running

**Visual indicators not showing:**
- Try refreshing the page
- Check if CSS is loading properly
- Ensure JavaScript is enabled