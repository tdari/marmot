export interface ConfigField {
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

export interface ConnectionTypeMeta {
	id: string;
	name: string;
	description?: string;
	icon?: string;
	category?: string;
	config_spec: ConfigField[];
}

export interface Connection {
	id: string;
	name: string;
	type: string;
	description?: string;
	config: Record<string, any>; // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	tags?: string[];
	created_by: string;
	created_at: string;
	updated_at: string;
}

export interface CreateConnectionInput {
	name: string;
	type: string;
	description?: string;
	config: Record<string, any>; // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	tags?: string[];
}

export interface UpdateConnectionInput {
	name?: string;
	description?: string;
	config?: Record<string, any>; // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	tags?: string[];
}

export interface ListConnectionsResponse {
	connections: Connection[];
	total: number;
	limit: number;
	offset: number;
}
