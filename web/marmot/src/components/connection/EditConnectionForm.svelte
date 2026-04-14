<script lang="ts">
	import { fetchApi } from '$lib/api';
	import { toasts } from '$lib/stores/toast';
	import { Database, Tag } from 'lucide-svelte';
	import IconifyIcon from '@iconify/svelte';
	import type { Connection } from '$lib/connections/types';
	import DynamicConfigForm from './DynamicConfigForm.svelte';
	import { onMount } from 'svelte';

	interface ConfigField {
		name: string;
		label: string;
		description?: string;
		type: string;
		required?: boolean;
		default?: any;
		sensitive?: boolean;
		placeholder?: string;
		options?: { value: string; label: string }[];
		fields?: ConfigField[];
		is_array?: boolean;
		show_when?: { field: string; value: any };
	}

	interface ConnectionTypeMeta {
		id: string;
		name: string;
		description?: string;
		icon?: string;
		category?: string;
		config_spec: ConfigField[];
	}

	let { connection, onCancel, onUpdate } = $props<{
		connection: Connection;
		onCancel: () => void;
		onUpdate: (updatedConnection: Connection) => void;
	}>();

	let loading = $state(false);
	let loadingTypes = $state(false);
	let connectionTypes: ConnectionTypeMeta[] = $state([]);
	let selectedType: ConnectionTypeMeta | null = $state(null);
	let fieldErrors: Record<string, string> = $state({});
	let error = $state<string | null>(null);

	function mergeConfigWithSpec(
		existingConfig: Record<string, any>,
		spec: ConfigField[]
	): Record<string, any> {
		const merged: Record<string, any> = {};

		spec.forEach((field) => {
			if (field.type === 'object' && field.fields) {
				if (field.is_array) {
					// Keep existing array or initialize as empty
					merged[field.name] = existingConfig[field.name] || [];
				} else {
					// Recursively merge nested objects
					merged[field.name] = mergeConfigWithSpec(
						existingConfig[field.name] || {},
						field.fields
					);
				}
			} else if (field.type === 'multiselect') {
				// Keep existing array or initialize as empty
				merged[field.name] = existingConfig[field.name] || [];
			} else {
				// Keep existing value or use default
				merged[field.name] =
					existingConfig[field.name] !== undefined ? existingConfig[field.name] : field.default;
			}
		});

		return merged;
	}

	// Initialize with current values
	let editedConnection = {
		name: connection.name,
		description: connection.description || '',
		config: { ...connection.config } as Record<string, any>,
		tags: connection.tags?.join(', ') || ''
	};

	// Form validation
	let isFormValid = $derived(editedConnection.name.trim() !== '');

	async function fetchConnectionTypes() {
		try {
			loadingTypes = true;
			const response = await fetchApi('/connections/types');
			if (!response.ok) {
				throw new Error('Failed to fetch connection types');
			}
			connectionTypes = await response.json();

			// Find the type meta for this connection
			selectedType = connectionTypes.find((t) => t.id === connection.type) || null;

			// Merge existing config with spec to ensure nested objects are properly initialized
			if (selectedType?.config_spec) {
				editedConnection.config = mergeConfigWithSpec(editedConnection.config, selectedType.config_spec);
			}
		} catch (err) {
			console.error('Failed to fetch connection types:', err);
			toasts.error('Failed to load connection types');
		} finally {
			loadingTypes = false;
		}
	}

	function clearFieldError(fieldPath: string) {
		delete fieldErrors[fieldPath];
		fieldErrors = { ...fieldErrors };
	}

	async function updateConnection() {
		if (!editedConnection.name.trim()) {
			error = 'Connection name is required';
			return;
		}

		try {
			loading = true;
			error = null;
			fieldErrors = {};

			const tags = editedConnection.tags
				.split(',')
				.map((t: string) => t.trim())
				.filter((t: string) => t);

			const response = await fetchApi(`/connections/${connection.id}`, {
				method: 'PUT',
				body: JSON.stringify({
					name: editedConnection.name !== connection.name ? editedConnection.name : undefined,
					description:
						editedConnection.description !== connection.description
							? editedConnection.description
							: undefined,
					config: editedConnection.config,
					tags: tags.length > 0 ? tags : []
				})
			});

			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error || 'Failed to update connection');
			}

			const updatedConnection = await response.json();
			toasts.success('Connection updated successfully');
			onUpdate(updatedConnection);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update connection';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchConnectionTypes();
	});
</script>

<div
	class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 m-4"
>
	<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-6 flex items-center">
		<Database class="h-5 w-5 mr-2 text-gray-500 dark:text-gray-400" />
		Edit Connection
	</h3>

	{#if error}
		<div
			class="mb-6 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/50 rounded-lg p-4"
		>
			<div class="flex items-start">
				<IconifyIcon
					icon="material-symbols:error"
					class="h-5 w-5 text-red-400 mt-0.5 flex-shrink-0"
				/>
				<p class="ml-3 text-sm text-red-700 dark:text-red-300">{error}</p>
			</div>
		</div>
	{/if}

	{#if loadingTypes}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-earthy-terracotta-600"></div>
		</div>
	{:else}
		<div class="space-y-6">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<label for="edit-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
						Connection Name
						<span class="text-red-500">*</span>
					</label>
					<input
						id="edit-name"
						type="text"
						bind:value={editedConnection.name}
						required
						class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
					/>
				</div>

				<div>
					<p class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
						Type
					</p>
					<div
						class="px-3 py-2 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-md text-sm text-gray-500 dark:text-gray-400"
					>
						{selectedType?.name || connection.type}
					</div>
				</div>
			</div>

			<div>
				<label for="edit-description" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
					Description (optional)
				</label>
				<input
					id="edit-description"
					type="text"
					bind:value={editedConnection.description}
					placeholder="Production database"
					class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
				/>
			</div>

			<div>
				<label
					for="edit-tags"
					class="flex text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 items-center"
				>
					<Tag class="h-4 w-4 mr-2 text-gray-500 dark:text-gray-400" />
					Tags (comma-separated, optional)
				</label>
				<input
					id="edit-tags"
					type="text"
					bind:value={editedConnection.tags}
					placeholder="production, analytics"
					class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
				/>
			</div>

			{#if selectedType?.config_spec}
				<div class="border-t border-gray-200 dark:border-gray-700 pt-6">
					<h4 class="text-sm font-medium text-gray-900 dark:text-gray-100 mb-4">Configuration</h4>
					<DynamicConfigForm
						configSpec={selectedType.config_spec}
						bind:config={editedConnection.config}
						{fieldErrors}
						onFieldChange={clearFieldError}
					/>
				</div>
			{/if}

			<div class="border-t border-gray-200 dark:border-gray-700 pt-6">
				<div class="text-xs text-gray-500 dark:text-gray-400 space-y-1">
					<p><strong>Created by:</strong> {connection.created_by}</p>
					<p><strong>Created at:</strong> {new Date(connection.created_at).toLocaleString()}</p>
					<p><strong>Updated at:</strong> {new Date(connection.updated_at).toLocaleString()}</p>
				</div>
			</div>
		</div>

		<div class="flex justify-end space-x-3 mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
			<button
				type="button"
				class="px-4 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 text-sm font-medium"
				onclick={onCancel}
			>
				Cancel
			</button>
			<button
				type="button"
				class="px-4 py-2 bg-earthy-terracotta-600 text-white rounded-md hover:bg-earthy-terracotta-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
				onclick={updateConnection}
				disabled={loading || !isFormValid}
			>
				{#if loading}
					<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
				{/if}
				Save Changes
			</button>
		</div>
	{/if}
</div>
