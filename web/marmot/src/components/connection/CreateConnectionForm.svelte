<script lang="ts">
	import { fetchApi } from '$lib/api';
	import { Database } from 'lucide-svelte';
	import type { Connection } from '$lib/connections/types';
	import DynamicConfigForm from './DynamicConfigForm.svelte';
	import { onMount } from 'svelte';
	import IconifyIcon from '@iconify/svelte';

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

	let { 
		onConnectionCreated = () => {},
		onCancel = () => {} 
	}: {
		onConnectionCreated?: (connection: Connection) => void;
		onCancel?: () => void;
	} = $props();

	let loading = $state(false);
	let loadingTypes = $state(false);
	let error = $state<string | null>(null);
	let connectionTypes = $state<ConnectionTypeMeta[]>([]);
	let selectedType = $state<ConnectionTypeMeta | null>(null);
	let fieldErrors = $state<Record<string, string>>({});

	let newConnection = $state({
		name: '',
		type: 'postgresql',
		description: '',
		config: {} as Record<string, any>,
		tags: ''
	});

	// Simple form validation - only check name is filled
	// Backend will validate the config structure based on ConfigSpec
	let isFormValid = $derived(newConnection.name.trim() !== '');

	function initializeConfigFromSpec(spec: ConfigField[]): Record<string, any> {
		const config: Record<string, any> = {};

		spec.forEach((field) => {
			if (field.type === 'object' && field.fields) {
				if (field.is_array) {
					// Initialize arrays of objects as empty arrays
					config[field.name] = [];
				} else {
					// Recursively initialize nested objects
					config[field.name] = initializeConfigFromSpec(field.fields);
				}
			} else if (field.type === 'multiselect') {
				// Initialize multiselect as empty array
				config[field.name] = [];
			} else if (field.default !== undefined) {
				// Use default values for other fields
				config[field.name] = field.default;
			}
		});

		return config;
	}

	async function fetchConnectionTypes() {
		try {
			loadingTypes = true;
			const response = await fetchApi('/connections/types');
			if (!response.ok) {
				throw new Error('Failed to fetch connection types');
			}
			const types = await response.json();
			connectionTypes = types.sort((a: ConnectionTypeMeta, b: ConnectionTypeMeta) =>
				a.name.localeCompare(b.name)
			);
			if (connectionTypes.length > 0) {
				handleTypeChange(connectionTypes[0].id);
			}
		} catch (err) {
			console.error('Failed to fetch connection types:', err);
			error = 'Failed to load connection types';
		} finally {
			loadingTypes = false;
		}
	}

	function handleTypeChange(typeId: string) {
		newConnection.type = typeId;
		selectedType = connectionTypes.find((t) => t.id === typeId) || null;

		// Initialize config with defaults from ConfigSpec, including nested objects
		if (selectedType?.config_spec) {
			newConnection.config = initializeConfigFromSpec(selectedType.config_spec);
		}

		fieldErrors = {};
	}

	function clearFieldError(fieldPath: string) {
		delete fieldErrors[fieldPath];
		fieldErrors = { ...fieldErrors };
	}

	async function createConnection() {
		try {
			loading = true;
			error = null;
			fieldErrors = {};

			const tags = newConnection.tags
				.split(',')
				.map((t) => t.trim())
				.filter((t) => t);

			const response = await fetchApi('/connections', {
				method: 'POST',
				body: JSON.stringify({
					name: newConnection.name,
					type: newConnection.type,
					description: newConnection.description || undefined,
					config: newConnection.config,
					tags: tags.length > 0 ? tags : undefined
				})
			});

			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error || 'Failed to create connection');
			}

			const createdConnection = await response.json();
			onConnectionCreated(createdConnection);

			// Reset form
			newConnection = {
				name: '',
				type: connectionTypes[0]?.id || 'postgresql',
				description: '',
				config: {},
				tags: ''
			};
			if (connectionTypes.length > 0) {
				handleTypeChange(connectionTypes[0].id);
			}
			error = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create connection';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchConnectionTypes();
	});
</script>

<div
	class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 m-4 animate-slide-down"
>
	<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-6 flex items-center">
		<Database class="h-5 w-5 mr-2 text-gray-500 dark:text-gray-400" />
		Create New Connection
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
					<label
						for="name"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
					>
						Connection Name
						<span class="text-red-500">*</span>
					</label>
					<input
						type="text"
						id="name"
						bind:value={newConnection.name}
						placeholder="my-database"
						required
						class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
					/>
				</div>

				<div>
					<label for="type" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
						Type
						<span class="text-red-500">*</span>
					</label>
					<select
						id="type"
						value={newConnection.type}
						onchange={(e) => handleTypeChange(e.currentTarget.value)}
						class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
					>
						{#each connectionTypes as type}
							<option value={type.id}>{type.name}</option>
						{/each}
					</select>
				</div>
			</div>

			<div>
				<label
					for="description"
					class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
				>
					Description (optional)
				</label>
				<input
					type="text"
					id="description"
					bind:value={newConnection.description}
					placeholder="Production database"
					class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
				/>
			</div>

			<div>
				<label for="tags" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
					Tags (comma-separated, optional)
				</label>
				<input
					type="text"
					id="tags"
					bind:value={newConnection.tags}
					placeholder="production, analytics"
					class="w-full px-3 py-2 bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-500 dark:focus:ring-earthy-terracotta-500 focus:border-transparent"
				/>
			</div>

			{#if selectedType?.config_spec}
				 <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
					<h4 class="text-sm font-medium text-gray-900 dark:text-gray-100 mb-4">
						Configuration
					</h4>
					<DynamicConfigForm
						configSpec={selectedType.config_spec}
						bind:config={newConnection.config}
						{fieldErrors}
						onFieldChange={clearFieldError}
					/>
				</div>
			{/if}

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
				onclick={createConnection}
				disabled={loading || !isFormValid}
			>
				{#if loading}
					<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
				{/if}
				Create Connection
			</button>
		</div>
	{/if}
</div>

<style>
	.animate-slide-down {
		animation: slideDown 0.2s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
