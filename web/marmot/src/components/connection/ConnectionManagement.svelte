<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchApi } from '$lib/api';
	import CreateConnectionForm from './CreateConnectionForm.svelte';
	import ConnectionTable from './ConnectionTable.svelte';
	import type { Connection } from '$lib/connections/types';

	let connections: Connection[] = [];
	let totalConnections = 0;
	let offset = 0;
	let limit = 10;
	let searchQuery = '';
	let creatingConnection = false;
	let editingConnectionId: string | null = null;
	let loading = false;
	let error: string | null = null;
	let searchTimer: ReturnType<typeof setTimeout>;

	async function handleConnectionCreated() {
		creatingConnection = false;
		await fetchConnections();
	}

	async function fetchConnections() {
		try {
			loading = true;
			const params = new URLSearchParams({
				limit: limit.toString(),
				offset: offset.toString(),
				...(searchQuery && { query: searchQuery })
			});

			const response = await fetchApi(`/connections?${params}`);
			const data = await response.json();
			connections = data.connections || [];
			totalConnections = data.total || 0;
		} catch (err) {
			error = err instanceof Error ? err.message : 'An error occurred';
		} finally {
			loading = false;
		}
	}

	onMount(fetchConnections);

	$: {
		if (searchQuery !== undefined) {
			if (searchTimer) clearTimeout(searchTimer);
			searchTimer = setTimeout(() => {
				offset = 0;
				fetchConnections();
			}, 300);
		}
	}

	async function handleConnectionUpdated(updatedConnection: Connection) {
		connections = connections.map((c) => (c.id === updatedConnection.id ? updatedConnection : c));
		editingConnectionId = null;
		await fetchConnections();
	}

	async function handleConnectionDeleted(connectionId: string) {
		connections = connections.filter((c) => c.id !== connectionId);
		await fetchConnections();
	}
</script>

<div
	class="bg-earthy-brown-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700"
>
	<div class="p-6">
		<div class="flex justify-between items-center mb-6">
			<div class="flex-1 max-w-md">
				<input
					type="text"
					placeholder="Search connections..."
					bind:value={searchQuery}
					class="w-full px-4 py-2 rounded-md border border-gray-300 dark:border-gray-600 focus:ring-2 focus:ring-earthy-terracotta-600 dark:focus:ring-earthy-terracotta-600 focus:border-transparent"
				/>
			</div>
			<button
				class="ml-4 px-4 py-2 bg-earthy-terracotta-700 dark:bg-earthy-terracotta-700 text-white rounded-md hover:bg-earthy-terracotta-800 dark:hover:bg-earthy-terracotta-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-earthy-terracotta-600 dark:focus:ring-earthy-terracotta-600"
				on:click={() => (creatingConnection = !creatingConnection)}
			>
				{creatingConnection ? 'Cancel' : 'Add Connection'}
			</button>
		</div>

		{#if creatingConnection}
			<CreateConnectionForm
				onConnectionCreated={handleConnectionCreated}
				onCancel={() => (creatingConnection = false)}
			/>
		{/if}

		{#if loading && !connections.length}
			<div class="flex justify-center p-8">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-earthy-terracotta-700"></div>
			</div>
		{:else if error}
			<div class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
				{error}
			</div>
		{:else}
			<ConnectionTable
				{connections}
				{editingConnectionId}
				onEdit={(connectionId) => (editingConnectionId = connectionId)}
				onUpdate={handleConnectionUpdated}
				onDelete={handleConnectionDeleted}
			/>

			<div class="mt-4 flex items-center justify-between">
				<div class="flex-1 flex justify-between items-center">
					<p class="text-sm text-gray-700 dark:text-gray-300">
						Showing {offset + 1} to {Math.min(offset + connections.length, totalConnections)} of {totalConnections}
						connections
					</p>
					<div class="flex space-x-2">
						<button
							class="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded-md text-sm disabled:opacity-50"
							disabled={offset === 0}
							on:click={() => { offset = Math.max(0, offset - limit); fetchConnections(); }}
						>
							Previous
						</button>
						<button
							class="px-3 py-1 border border-gray-300 dark:border-gray-600 rounded-md text-sm disabled:opacity-50"
							disabled={offset + connections.length >= totalConnections}
							on:click={() => { offset = offset + limit; fetchConnections(); }}
						>
							Next
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
