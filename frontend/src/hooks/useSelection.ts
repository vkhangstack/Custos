import { useState, useCallback } from 'react';

export interface UseSelectionResult<T> {
    selectedIds: Set<T>;
    isSelected: (id: T) => boolean;
    toggle: (id: T) => void;
    toggleAll: (ids: T[]) => void;
    selectAll: (ids: T[]) => void;
    clear: () => void;
    count: number;
    allSelected: (ids: T[]) => boolean;
    isPartiallySelected: (ids: T[]) => boolean;
}

export function useSelection<T = string>(initialSelected: T[] = []): UseSelectionResult<T> {
    const [selectedIds, setSelectedIds] = useState<Set<T>>(new Set(initialSelected));

    const isSelected = useCallback((id: T) => selectedIds.has(id), [selectedIds]);

    const toggle = useCallback((id: T) => {
        setSelectedIds(prev => {
            const next = new Set(prev);
            if (next.has(id)) {
                next.delete(id);
            } else {
                next.add(id);
            }
            return next;
        });
    }, []);

    const selectAll = useCallback((ids: T[]) => {
        setSelectedIds(new Set(ids));
    }, []);

    const clear = useCallback(() => {
        setSelectedIds(new Set());
    }, []);

    const toggleAll = useCallback((ids: T[]) => {
        setSelectedIds(prev => {
            // If all are currently selected, clear them. Otherwise select all.
            const allCurrentlySelected = ids.every(id => prev.has(id));

            if (allCurrentlySelected) {
                return new Set();
            } else {
                return new Set(ids);
            }
        });
    }, []);

    const allSelected = useCallback((ids: T[]) => {
        return ids.length > 0 && ids.every(id => selectedIds.has(id));
    }, [selectedIds]);

    const isPartiallySelected = useCallback((ids: T[]) => {
        const selectedCount = ids.filter(id => selectedIds.has(id)).length;
        return selectedCount > 0 && selectedCount < ids.length;
    }, [selectedIds]);

    return {
        selectedIds,
        isSelected,
        toggle,
        toggleAll,
        selectAll,
        clear,
        count: selectedIds.size,
        allSelected,
        isPartiallySelected
    };
}
