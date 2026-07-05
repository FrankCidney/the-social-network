const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? '';

type RequestOptions = Omit<RequestInit, 'body'> & {
  body?: BodyInit | Record<string, unknown> | unknown[];
};

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers);
  const isFormData = options.body instanceof FormData;

  if (options.body !== undefined && !isFormData && !headers.has('content-type')) {
    headers.set('content-type', 'application/json');
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
    credentials: 'include',
    body:
      options.body === undefined
        ? undefined
        : isFormData
          ? options.body
          : JSON.stringify(options.body),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

async function getErrorMessage(response: Response) {
	try {
		const data = await response.json();

		if (typeof data?.error?.message === 'string') {
			return data.error.message;
		}

		if (typeof data?.error === 'string') {
			return data.error;
		}

		if (typeof data?.message === 'string') {
			return data.message;
		}
	} catch {
		// Fall back to the status text below when the response is not JSON.
	}

	return response.statusText || 'Request failed';
}

export type RegisterPayload = {
	email: string;
	password: string;
	first_name: string;
	last_name: string;
	dob: string;
	nickname?: string;
};

export type LoginPayload = {
	email: string;
	password: string;
};

export type AuthResponse = {
	user: PublicUser;
	token: string;
};

export type PublicUser = {
	id: string;
	first_name: string;
	last_name: string;
	nickname?: string;
	avatar_path?: string;
	is_public: boolean;
};

export type PostResponse = {
	id: string;
	author: PublicUser;
	group_id?: string;
	content?: string;
	image_url?: string;
	privacy: string;
	created_at: string;
};

export type PostListResponse = {
	posts: PostResponse[];
	total: number;
	limit: number;
	offset: number;
};

export type CreatePostPayload = {
	content?: string;
	privacy: string;
	visible_to?: string[];
	group_id?: string;
};

export type UploadPostImageResponse = {
	image_url: string;
};

export type CreatedPostResponse = {
	id: string;
	group_id?: string;
	content?: string;
	image_url?: string;
	privacy: string;
	created_at: string;
};

export type Group = {
	id: string;
	creator_id: string;
	title: string;
	description: string;
	created_at: string;
};

export const authAPI = {
	register(payload: RegisterPayload) {
		return request<AuthResponse>('/api/auth/register', {
			method: 'POST',
			body: payload,
		});
	},

	login(payload: LoginPayload) {
		return request<AuthResponse>('/api/auth/login', {
			method: 'POST',
			body: payload,
		});
	},
};

export const groupAPI = {
	getGroups() {
		return request<Group[] | { groups?: Group[] }>('/api/groups');
	},
};

export const feedAPI = {
	getFeed(limit?: number, offset?: number) {
		const params = new URLSearchParams();
		if (limit !== undefined) params.set('limit', String(limit));
		if (offset !== undefined) params.set('offset', String(offset));
		const query = params.toString();
		return request<PostListResponse>(`/api/feed${query ? `?${query}` : ''}`);
	},

	createPost(payload: CreatePostPayload) {
		return request<CreatedPostResponse>('/api/posts', {
			method: 'POST',
			body: payload,
		});
	},

	getPost(postId: string) {
		return request<PostResponse>(`/api/posts/${postId}`);
	},

	uploadPostImage(postId: string, file: File) {
		const formData = new FormData();
		formData.append('image', file);

		return request<UploadPostImageResponse>(`/api/posts/${postId}/image`, {
			method: 'POST',
			body: formData,
		});
	},
};

export function resolveAssetUrl(path?: string) {
	if (!path) {
		return undefined;
	}

	if (/^https?:\/\//i.test(path)) {
		return path;
	}

	const normalizedPath = path.startsWith('/') ? path : `/${path}`;
	return API_BASE_URL ? `${API_BASE_URL}${normalizedPath}` : normalizedPath;
}
