<script lang="ts">
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

	interface Props {
		configSpec: ConfigField[];
		config: Record<string, any>;
		fieldErrors?: Record<string, string>;
		onFieldChange?: (fieldPath: string) => void;
	}

	let { configSpec, config = $bindable(), fieldErrors = {}, onFieldChange }: Props = $props();

	let expandedSections: Record<string, boolean> = $state({});

	function getFieldType(field: ConfigField): string {
		if (field.type === 'bool' || field.type === 'boolean') return 'checkbox';
		if (field.type === 'int' || field.type === 'number' || field.type === 'integer')
			return 'number';
		if (field.type === 'password' || field.sensitive) return 'password';
		if (field.name.toLowerCase() === 'url') return 'url';
		return 'text';
	}

	function toggleSection(sectionName: string) {
		expandedSections[sectionName] = !expandedSections[sectionName];
	}

	function isExpanded(sectionName: string): boolean {
		if (expandedSections[sectionName] !== undefined) {
			return expandedSections[sectionName];
		}
		return true; // Default to expanded
	}

	function shouldHideField(
		field: ConfigField,
		configObj: Record<string, any>,
		rootConfig?: Record<string, any>
	): boolean {
		if (configObj.use_default === true && field.name !== 'use_default') {
			return true;
		}
		if (field.show_when) {
			const checkObj = rootConfig || configObj;
			const currentValue = checkObj[field.show_when.field];
			if (currentValue !== field.show_when.value) {
				return true;
			}
		}
		return false;
	}

	function clearFieldError(fieldPath: string) {
		if (onFieldChange) {
			onFieldChange(fieldPath);
		}
	}
</script>

{#snippet renderField(
	field: ConfigField,
	fieldPath: string,
	configObj: Record<string, any>,
	depth: number = 0
)}
	{#if field.type === 'object' && field.is_array && field.fields}
		<!-- Array of objects -->
		<div class="md:col-span-2">
			<div class="block">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">
					{field.label}
					{#if field.required}
						<span class="text-red-500">*</span>
					{/if}
				</span>
				{#if field.description}
					<p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
						{field.description}
					</p>
				{/if}
				{#if true}
					{@const arrayValue = (configObj[field.name] ??= [])}
					<div class="space-y-3">
						{#each arrayValue as item, index}
							<div
								class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-gray-50/50 dark:bg-gray-750/50"
							>
								<div class="flex items-end justify-end mb-3">
									<button
										type="button"
										onclick={(e) => {
											e.preventDefault();
											arrayValue.splice(index, 1);
											configObj[field.name] = [...arrayValue];
											clearFieldError(fieldPath);
										}}
										class="p-1 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
									>
										<IconifyIcon icon="material-symbols:close" class="h-4 w-4" />
									</button>
								</div>
								<div class="grid grid-cols-1 gap-3">
									{#each field.fields as nestedField}
										{@const nestedItemPath = `${fieldPath}[${index}].${nestedField.name}`}
										<div>
											<label class="block">
												<span
													class="text-xs font-medium text-gray-700 dark:text-gray-300 mb-1 block"
												>
													{nestedField.label}
													{#if nestedField.required}
														<span class="text-red-500">*</span>
													{/if}
												</span>
												<input
													type={getFieldType(nestedField)}
													bind:value={item[nestedField.name]}
													oninput={(e) => {
														configObj[field.name] = [...arrayValue];
														clearFieldError(nestedItemPath);
													}}
													placeholder={nestedField.placeholder}
													required={nestedField.required}
													data-field-path={nestedItemPath}
													class="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-600 focus:border-transparent transition-all {fieldErrors[
														nestedItemPath
													]
														? 'border-red-500 dark:border-red-500'
														: ''}"
												/>
												{#if fieldErrors[nestedItemPath]}
													<p
														class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-start"
													>
														<IconifyIcon
															icon="material-symbols:error"
															class="h-4 w-4 mr-1 mt-0.5 flex-shrink-0"
														/>
														{fieldErrors[nestedItemPath]}
													</p>
												{/if}
											</label>
										</div>
									{/each}
								</div>
							</div>
						{/each}
						<button
							type="button"
							onclick={(e) => {
								e.preventDefault();
								const newItem: Record<string, any> = {};
								field.fields?.forEach((f) => {
									newItem[f.name] = f.default || '';
								});
								arrayValue.push(newItem);
								configObj[field.name] = [...arrayValue];
								clearFieldError(fieldPath);
							}}
							class="w-full px-4 py-2.5 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-600 dark:text-gray-400 hover:border-earthy-terracotta-600 hover:text-earthy-terracotta-600 dark:hover:text-earthy-terracotta-400 transition-colors flex items-center justify-center gap-2"
						>
							<IconifyIcon icon="material-symbols:add" class="h-5 w-5" />
							Add {field.label}
						</button>
					</div>
				{/if}
				{#if fieldErrors[fieldPath]}
					<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-start">
						<IconifyIcon
							icon="material-symbols:error"
							class="h-4 w-4 mr-1 mt-0.5 flex-shrink-0"
						/>
						{fieldErrors[fieldPath]}
					</p>
				{/if}
			</div>
		</div>
	{:else if field.type === 'object' && field.fields}
		<div class="md:col-span-2">
			<div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
				<button
					type="button"
					onclick={() => toggleSection(fieldPath)}
					class="w-full flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors text-left"
				>
					<div class="flex items-center">
						<IconifyIcon
							icon={isExpanded(fieldPath)
								? 'material-symbols:expand-more'
								: 'material-symbols:chevron-right'}
							class="h-5 w-5 text-gray-500 dark:text-gray-400 transition-transform"
						/>
						<span class="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300">
							{field.label}
							{#if field.required}
								<span class="text-red-500 ml-1">*</span>
							{/if}
						</span>
					</div>
					{#if field.description}
						<span class="text-xs text-gray-500 dark:text-gray-400 ml-2 truncate"
							>{field.description}</span
						>
					{/if}
				</button>
				{#if isExpanded(fieldPath)}
					<div
						class="px-4 pb-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-750/50"
					>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4">
							{#each field.fields as nestedField}
								{@const nestedPath = `${fieldPath}.${nestedField.name}`}
								{@const nestedConfigObj = (configObj[field.name] ??= {})}
								{#if !shouldHideField(nestedField, nestedConfigObj, config)}
									{@render renderField(nestedField, nestedPath, nestedConfigObj, depth + 1)}
								{/if}
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{:else if field.type === 'bool' || field.type === 'boolean'}
		<div class="md:col-span-2">
			<label
				class="flex items-start p-3 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-750 cursor-pointer transition-colors {fieldErrors[
					fieldPath
				]
					? 'border-red-500 dark:border-red-500'
					: 'border-gray-200 dark:border-gray-700'}"
				data-field-path={fieldPath}
			>
				<input
					type="checkbox"
					bind:checked={configObj[field.name]}
					onchange={() => clearFieldError(fieldPath)}
					class="h-4 w-4 mt-0.5 text-earthy-terracotta-700 focus:ring-earthy-terracotta-600 border-gray-300 rounded"
				/>
				<div class="ml-3 flex-1">
					<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
						{field.label}
					</span>
					{#if field.description}
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
							{field.description}
						</p>
					{/if}
					{#if fieldErrors[fieldPath]}
						<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-start">
							<IconifyIcon
								icon="material-symbols:error"
								class="h-4 w-4 mr-1 mt-0.5 flex-shrink-0"
							/>
							{fieldErrors[fieldPath]}
						</p>
					{/if}
				</div>
			</label>
		</div>
	{:else if field.type === 'multiselect'}
		<!-- Array/List field -->
		<div class="md:col-span-2">
			<div class="block" data-field-path={fieldPath}>
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">
					{field.label}
					{#if field.required}
						<span class="text-red-500">*</span>
					{/if}
				</span>
				{#if field.description}
					<p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
						{field.description}
					</p>
				{/if}
				{#if true}
					{@const arrayValue = (configObj[field.name] ??= [])}
					<div class="space-y-2">
						{#each arrayValue as item, index}
							<div class="flex items-center gap-2">
								<input
									type="text"
									value={item}
									oninput={(e) => {
										const target = e.target as HTMLInputElement;
										arrayValue[index] = target.value;
										configObj[field.name] = [...arrayValue];
										clearFieldError(fieldPath);
									}}
									placeholder={field.placeholder}
									class="flex-1 px-4 py-2.5 border rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-600 focus:border-transparent transition-all {fieldErrors[
										fieldPath
									]
										? 'border-red-500 dark:border-red-500'
										: 'border-gray-300 dark:border-gray-600'}"
								/>
								<button
									type="button"
									onclick={(e) => {
										e.preventDefault();
										e.stopPropagation();
										arrayValue.splice(index, 1);
										configObj[field.name] = [...arrayValue];
										clearFieldError(fieldPath);
									}}
									class="p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
									aria-label="Remove item"
								>
									<IconifyIcon icon="material-symbols:close" class="h-5 w-5" />
								</button>
							</div>
						{/each}
						<div class="flex items-center gap-2">
							<input
								type="text"
								placeholder={`Type to add ${field.label.toLowerCase()}...`}
								onkeydown={(e) => {
									if (e.key === 'Enter') {
										e.preventDefault();
										const target = e.target as HTMLInputElement;
										const value = target.value.trim();
										if (value) {
											const newArray = [...arrayValue, value];
											configObj[field.name] = newArray;
											target.value = '';
											clearFieldError(fieldPath);
										}
									}
								}}
								class="flex-1 px-4 py-2.5 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-600 focus:border-earthy-terracotta-600 transition-all"
							/>
						</div>
						<p class="text-xs text-gray-500 dark:text-gray-400">Press Enter to add items</p>
					</div>
				{/if}
				{#if fieldErrors[fieldPath]}
					<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-start">
						<IconifyIcon
							icon="material-symbols:error"
							class="h-4 w-4 mr-1 mt-0.5 flex-shrink-0"
						/>
						{fieldErrors[fieldPath]}
					</p>
				{/if}
			</div>
		</div>
	{:else}
		<div>
			<label class="block">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">
					{field.label}
					{#if field.required}
						<span class="text-red-500">*</span>
					{/if}
				</span>
				{#if field.description}
					<p class="text-xs text-gray-500 dark:text-gray-400 mb-2">
						{field.description}
					</p>
				{/if}
				{#if field.options && field.options.length > 0}
					<select
						value={configObj[field.name] ?? field.default ?? ''}
						onchange={(e) => {
							configObj[field.name] = (e.target as HTMLSelectElement).value;
							clearFieldError(fieldPath);
						}}
						data-field-path={fieldPath}
						class="w-full px-4 py-2.5 border rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-600 focus:border-transparent transition-all {fieldErrors[
							fieldPath
						]
							? 'border-red-500 dark:border-red-500'
							: 'border-gray-300 dark:border-gray-600'}"
						required={field.required}
					>
						<option value="">Select...</option>
						{#each field.options as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				{:else}
					<input
						type={getFieldType(field)}
						value={configObj[field.name] ?? field.default ?? ''}
						oninput={(e) => {
							const target = e.target as HTMLInputElement;
							configObj[field.name] =
								field.type === 'int' || field.type === 'number'
									? Number(target.value)
									: target.value;
							clearFieldError(fieldPath);
						}}
						placeholder={field.placeholder || (field.default ? String(field.default) : '')}
						data-field-path={fieldPath}
						class="w-full px-4 py-2.5 border rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-earthy-terracotta-600 focus:border-transparent transition-all {field.type ===
						'password'
							? 'font-mono'
							: ''} {fieldErrors[fieldPath]
							? 'border-red-500 dark:border-red-500'
							: 'border-gray-300 dark:border-gray-600'}"
						required={field.required}
					/>
				{/if}
				{#if fieldErrors[fieldPath]}
					<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-start">
						<IconifyIcon
							icon="material-symbols:error"
							class="h-4 w-4 mr-1 mt-0.5 flex-shrink-0"
						/>
						{fieldErrors[fieldPath]}
					</p>
				{/if}
			</label>
		</div>
	{/if}
{/snippet}

<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
	{#each configSpec as field}
		{#if !shouldHideField(field, config)}
			{@render renderField(field, field.name, config, 0)}
		{/if}
	{/each}
</div>
