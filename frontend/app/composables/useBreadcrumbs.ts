export const useBreadcrumbs = () => {
  const breadcrumbs = useState<{ label: string; to?: string }[]>('breadcrumbs', () => [])

  const setBreadcrumbs = (items: { label: string; to?: string }[]) => {
    breadcrumbs.value = items
  }

  const clearBreadcrumbs = () => {
    breadcrumbs.value = []
  }

  return {
    breadcrumbs,
    setBreadcrumbs,
    clearBreadcrumbs,
  }
}
