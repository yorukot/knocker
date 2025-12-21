<script lang="ts">
  import Icon from '@iconify/svelte';
  import * as Card from '$lib/components/ui/card/index.js';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator/index.js';
  import type { Incident } from '../../../types';
  import type { IncidentWithMonitors } from './+page';

  /** @type {import('./$types').PageProps} */
  let { data } = $props();

  const incidents = $derived<IncidentWithMonitors[]>(data.incidents ?? []);
  const openIncidents = $derived(sortIncidents(incidents.filter((i) => i.status !== 'resolved')));
  const resolvedIncidents = $derived(sortIncidents(incidents.filter((i) => i.status === 'resolved')));

  function sortIncidents<T extends Incident>(list: T[]): T[] {
    return [...list].sort((a, b) => {
      const aTime = new Date(a.resolvedAt ?? a.startedAt).getTime();
      const bTime = new Date(b.resolvedAt ?? b.startedAt).getTime();
      return bTime - aTime;
    });
  }

  function formatDate(value?: string) {
    if (!value) return '—';
    const d = new Date(value);
    return Number.isNaN(d.getTime()) ? '—' : d.toLocaleString();
  }

  function formatDuration(incident: Incident) {
    const start = new Date(incident.startedAt).getTime();
    const end = new Date(incident.resolvedAt ?? Date.now()).getTime();
    if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) return '—';
    const diffMs = end - start;
    const minutes = Math.floor(diffMs / 60000);
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours < 24) return `${hours}h ${mins}m`;
    const days = Math.floor(hours / 24);
    return `${days}d ${hours % 24}h`;
  }

  const severityMeta: Record<Incident['severity'], { label: string; badge: string; dot: string }> = {
    emergency: {
      label: 'Emergency',
      badge: '!bg-destructive !text-destructive-foreground border-transparent',
      dot: 'bg-destructive'
    },
    critical: {
      label: 'Critical',
      badge: '!bg-destructive/80 !text-destructive-foreground border-transparent',
      dot: 'bg-destructive/80'
    },
    major: {
      label: 'Major',
      badge: '!bg-amber-500 !text-amber-950 border-transparent',
      dot: 'bg-amber-500'
    },
    minor: {
      label: 'Minor',
      badge: '!bg-secondary !text-secondary-foreground border-transparent',
      dot: 'bg-secondary'
    },
    info: {
      label: 'Info',
      badge: '!bg-muted !text-foreground border-transparent',
      dot: 'bg-muted'
    }
  };

  const severityMetaSafe = (severity: Incident['severity']) =>
    severityMeta[severity] ?? {
      label: severity,
      badge: '!bg-muted !text-foreground border-transparent',
      dot: 'bg-muted'
    };

  const statusMeta: Record<Incident['status'], { label: string; color: string }> = {
    detected: { label: 'Detected', color: 'text-amber-600' },
    investigating: { label: 'Investigating', color: 'text-amber-700' },
    identified: { label: 'Identified', color: 'text-amber-700' },
    monitoring: { label: 'Monitoring', color: 'text-sky-700' },
    resolved: { label: 'Resolved', color: 'text-success' }
  };

  const dotClass = (severity: Incident['severity']) => severityMetaSafe(severity).dot;
  const badgeClass = (severity: Incident['severity']) => severityMetaSafe(severity).badge;
  const severityLabel = (severity: Incident['severity']) => severityMetaSafe(severity).label;

</script>

<div class="flex flex-col gap-6">
  <header class="flex flex-col gap-3">
    <div class="flex items-start justify-between gap-3 flex-wrap">
      <div class="flex flex-col gap-1">
        <p class="text-sm text-muted-foreground">Team incidents</p>
        <div class="flex items-center gap-2 flex-wrap">
          <h1 class="text-2xl font-bold">Incidents</h1>
          <Badge variant="secondary">Manual</Badge>
        </div>
      </div>
      <Button href="incidents/new">
        <Icon icon="lucide:plus" />
        Create incident
      </Button>
    </div>
    <p class="text-sm text-muted-foreground">
      Open incidents are shown first. Resolved ones appear below with their duration.
    </p>
  </header>

  <section class="flex flex-col gap-3">
    {@render SectionHeading({ label: 'Open', count: openIncidents.length, icon: 'lucide:alert-octagon' })}
    {#if openIncidents.length === 0}
      {@render EmptyState({ message: 'No open incidents' })}
    {:else}
      {@render IncidentList({ items: openIncidents })}
    {/if}
  </section>

  <Separator class="my-2" />

  <section class="flex flex-col gap-3">
    {@render SectionHeading({ label: 'Resolved', count: resolvedIncidents.length, icon: 'lucide:check-circle-2' })}
    {#if resolvedIncidents.length === 0}
      {@render EmptyState({ message: 'No resolved incidents yet' })}
    {:else}
      {@render IncidentList({ items: resolvedIncidents, subdued: true })}
    {/if}
  </section>
</div>

<!-- Components -->
{#snippet SectionHeading({ label, count, icon }: { label: string; count: number; icon: string })}
  <div class="flex items-center gap-2">
    <Icon icon={icon} class="size-5" />
    <h2 class="text-lg font-semibold">{label}</h2>
    <Badge variant="secondary">{count}</Badge>
  </div>
{/snippet}

{#snippet EmptyState({ message }: { message: string })}
  <Card.Root class="p-6 text-center text-muted-foreground border-dashed">
    {message}
  </Card.Root>
{/snippet}

{#snippet IncidentList({ items, subdued = false }: { items: IncidentWithMonitors[]; subdued?: boolean })}
  <div class="flex flex-col gap-3">
    {#each items as incident (incident.id)}
      {@render IncidentCard({ incident, subdued })}
    {/each}
  </div>
{/snippet}

{#snippet IncidentCard({ incident, subdued = false }: { incident: IncidentWithMonitors; subdued?: boolean })}
  <Card.Root class={`p-4 flex flex-col gap-3 ${subdued ? 'opacity-80' : ''}`}>
    <div class="flex items-start gap-3 flex-wrap">
      <div class="flex-1 min-w-0 flex flex-col gap-1">
        <div class="flex items-center gap-2 flex-wrap">
          <p class="text-lg font-semibold truncate">Incident #{incident.id}</p>
          <Badge>{severityMeta[incident.severity].label}</Badge>
          <span class={`text-sm font-medium ${statusMeta[incident.status].color}`}>
            {statusMeta[incident.status].label}
          </span>
          {#if incident.isPublic}
            <Badge variant="outline" class="gap-1">
              <Icon icon="lucide:globe-2" class="size-3.5" /> Public
            </Badge>
          {/if}
          {#if incident.autoResolve}
            <Badge variant="outline" class="gap-1">
              <Icon icon="lucide:clock-3" class="size-3.5" /> Auto-resolve
            </Badge>
          {/if}
        </div>
        <div class="text-sm text-muted-foreground flex gap-3 flex-wrap">
          <span class="flex items-center gap-1">
            <Icon icon="lucide:play-circle" class="size-4" />
            Started {formatDate(incident.startedAt)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:flag" class="size-4" />
            Resolved {formatDate(incident.resolvedAt)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:timer" class="size-4" />
            Duration {formatDuration(incident)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:clock" class="size-4" />
            Updated {formatDate(incident.updatedAt)}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Button size="sm" variant="ghost" href={`incidents/${incident.id}`}>
          View incident
          <Icon icon="lucide:arrow-right" />
        </Button>
      </div>
    </div>

    <Card.Content class="p-0 flex flex-col gap-2">
      <div class="flex items-center gap-2 text-sm font-medium">
        <Icon icon="lucide:monitor" class="size-4" />
        Monitors ({incident.monitorNames.length})
      </div>
      {#if incident.monitorNames.length === 0}
        <p class="text-sm text-muted-foreground">No monitor links found for this incident.</p>
      {:else}
        <div class="flex flex-wrap gap-2">
          {#each incident.monitorNames as name (name)}
            <Badge variant="secondary" class="truncate max-w-[14rem]">
              {name}
            </Badge>
          {/each}
        </div>
      {/if}
    </Card.Content>
  </Card.Root>
{/snippet}
