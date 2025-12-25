<script lang="ts">
  import Button from '$lib/components/ui/button/button.svelte';
  import * as Card from '$lib/components/ui/card';
  import * as InputGroup from '$lib/components/ui/input-group';
  import * as Select from '$lib/components/ui/select';
  import { Separator } from '$lib/components/ui/separator';
  import Icon from '@iconify/svelte';
  import type { StatusPageElement, StatusPageElementType } from '$lib/types';
  import StatusPageMonitorRow from './monitor-raw-status-page.svelte';

  const {
    group,
    namePrefix,
    onDeleteGroup,
    onDeleteMonitor
  }: {
    group: StatusPageElement;
    namePrefix: string;
    onDeleteGroup?: (groupId: string) => void;
    onDeleteMonitor?: (monitorId: string) => void;
  } = $props();

  let typeValue = $derived<StatusPageElementType>(group.type);

  const typeLabel = (t: StatusPageElementType) =>
    t === 'historical_timeline' ? 'Historical Timeline' : 'Only Current Status';

  const typeIcon = (t: StatusPageElementType) =>
    t === 'historical_timeline' ? 'lucide:chart-line' : 'lucide:circle-small';
</script>

<Card.Root class="bg-muted p-0">
  <Card.Content class="p-0">
    <div class="flex justify-between items-center p-2 gap-2">
      <InputGroup.Root class="w-full">
        <InputGroup.Input
          name={`${namePrefix}.name`}
          value={group.name}
          placeholder="Please enter element name"
        />
        <input type="hidden" name={`${namePrefix}.monitor`} value="false" />
        <input type="hidden" name={`${namePrefix}.sortOrder`} value={group.sortOrder} />
        <InputGroup.Addon class="hidden sm:block">
          <Icon icon="lucide:layers" />
        </InputGroup.Addon>
      </InputGroup.Root>

      <div class="flex items-center gap-2">
        <Select.Root type="single" bind:value={typeValue}>
          <Select.Trigger class="lg:w-51">
            <Icon icon={typeIcon(typeValue)} />
            <p class="hidden lg:block">{typeLabel(typeValue)}</p>
          </Select.Trigger>
          <Select.Content>
            <Select.Group>
              <Select.Item value="historical_timeline" label="Historical Timeline">
                <Icon icon="lucide:chart-line" /> Historical Timeline
              </Select.Item>
              <Select.Item value="current_status_indicator" label="Only Current Status">
                <Icon icon="lucide:circle-small" /> Only Current Status
              </Select.Item>
            </Select.Group>
          </Select.Content>
        </Select.Root>
        <input type="hidden" name={`${namePrefix}.type`} bind:value={typeValue} />

        <Button
          size="icon"
          variant="destructive"
          onclick={() => onDeleteGroup?.(group.id)}
        >
          <Icon icon="lucide:trash" />
        </Button>
      </div>
    </div>

    <Separator />

    {#if (group.monitors ?? []).length === 0}
      <div class="p-2 text-sm opacity-70">No monitor yet.</div>
    {:else}
      {#each group.monitors ?? [] as m, monitorIndex (m.id)}
        <StatusPageMonitorRow
          monitor={m}
          namePrefix={`${namePrefix}.monitors.${monitorIndex}`}
          onDelete={onDeleteMonitor}
        />
      {/each}
    {/if}
  </Card.Content>
</Card.Root>
