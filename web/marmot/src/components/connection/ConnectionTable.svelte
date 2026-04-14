<script lang="ts">
	import { fetchApi } from '$lib/api';
	import { toasts, handleApiError } from '$lib/stores/toast';
	import EditConnectionForm from './EditConnectionForm.svelte';
	import ConfirmModal from '$components/ui/ConfirmModal.svelte';
	import { Database, Tag } from 'lucide-svelte';
	import type { Connection } from '$lib/connections/types';

	export let connections: Connection[] = [];
	export let editingConnectionId: string | null = null;
	export let onEdit: (connectionId: string | null) => void;
	export let onUpdate: (connection: Connection) => void;
	export let onDelete: (connectionId: string) => void;

	let showDeleteModal = false;
	let connectionToDelete: Connection | null = null;
	let deleteSchedules = false;

	async function handleDelete() {
		if (!connectionToDelete) return;

		try {
			let url = `/connections/${connectionToDelete.id}`;
			if (deleteSchedules) {
				url += '?teardown=true';
			}

			const response = await fetchApi(url, {
				method: 'DELETE'
			});
			if (!response.ok) {
				const errorMsg = await handleApiError(response);
				toasts.error(errorMsg);
				return;
			}
			toasts.success(`Connection "${connectionToDelete.name}" deleted successfully`);
			onDelete(connectionToDelete.id);
			showDeleteModal = false;
			connectionToDelete = null;
			deleteSchedules = false;
		} catch (err) {
			toasts.error(err instanceof Error ? err.message : 'Failed to delete connection');
		}
	}

	function getTypeDisplay(type: string): string {
		const typeMap: Record<string, string> = {
			postgresql: 'PostgreSQL',
			mysql: 'MySQL',
			bigquery: 'BigQuery',
			s3: 'Amazon S3',
			snowflake: 'Snowflake',
			redshift: 'Redshift',
			databricks: 'Databricks'
		};
		return typeMap[type] || type.charAt(0).toUpperCase() + type.slice(1);
	}
</script>

<div class="overflow-x-auto">
	<table class="min-w-full">
		<thead>
			<tr>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-earthy-brown-100 dark:bg-gray-800"
					>Name</th
				>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-earthy-brown-100 dark:bg-gray-800"
					>Type</th
				>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-earthy-brown-100 dark:bg-gray-800"
					>Description</th
				>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider bg-earthy-brown-100 dark:bg-gray-800"
					>Tags</th
				>
				<th
					class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider bg-earthy-brown-100 dark:bg-gray-800"
					>Actions</th
				>
			</tr>
		</thead>
		<tbody class="divide-y divide-earthy-brown-100 bg-earthy-brown-50 dark:bg-gray-900">
			{#each connections as connection}
				<tr class="hover:bg-earthy-brown-100 dark:hover:bg-gray-800 transition-colors">
					{#if editingConnectionId === connection.id}
						<td colspan="5">
							<EditConnectionForm
								{connection}
								onCancel={() => onEdit(null)}
								onUpdate={(updatedConnection) => onUpdate(updatedConnection)}
							/>
						</td>
					{:else}
						<td
							class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-100"
						>
							<div class="flex items-center">
								<Database class="h-4 w-4 mr-2 text-gray-400" />
								{connection.name}
							</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
							>
								{getTypeDisplay(connection.type)}
							</span>
						</td>
						<td class="px-6 py-4 text-sm text-gray-600 dark:text-gray-400">
							{connection.description || '-'}
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							{#if connection.tags && connection.tags.length > 0}
								<div class="flex flex-wrap gap-1">
									{#each connection.tags as tag}
										<span
											class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-200"
										>
											<Tag class="h-3 w-3 mr-1" />
											{tag}
										</span>
									{/each}
								</div>
							{:else}
								<span class="text-sm text-gray-400">-</span>
							{/if}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
							<button
								type="button"
								class="text-earthy-terracotta-700 hover:text-earthy-terracotta-800 dark:text-earthy-terracotta-500 dark:hover:text-earthy-terracotta-400 mr-3"
								on:click={() => onEdit(connection.id)}
							>
								Edit
							</button>
							<button
								type="button"
								class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
								on:click={() => {
									connectionToDelete = connection;
									showDeleteModal = true;
								}}
							>
								Delete
							</button>
						</td>
					{/if}
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<ConfirmModal
	bind:show={showDeleteModal}
	title="Delete Connection"
	message="Are you sure you want to delete this connection? This action cannot be undone."
	confirmText="Delete"
	cancelText="Cancel"
	variant="danger"
	checkboxLabel="Delete all schedules using this connection"
	bind:checkboxChecked={deleteSchedules}
	onConfirm={handleDelete}
	onCancel={() => {
		showDeleteModal = false;
		connectionToDelete = null;
		deleteSchedules = false;
	}}
/>
