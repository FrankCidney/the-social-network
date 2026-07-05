'use client';

import * as React from 'react';
import { motion } from 'framer-motion';
import {
	Check,
	ChevronDown,
	FileText,
	Image as ImageIcon,
	LoaderCircle,
	MessageCircle,
	MoreHorizontal,
	Share2,
	X,
} from 'lucide-react';
import { Button } from '@/components/ui/Button';
import {
	feedAPI,
	FollowListResponse,
	PostListResponse,
	PostResponse,
	profileAPI,
	PublicUser,
	resolveAssetUrl,
} from '@/lib/api';

const DEFAULT_FEED_LIMIT = 20;
const FOLLOWERS_PAGE_SIZE = 20;

type PostPrivacy = 'public' | 'almost_private' | 'private';

const PRIVACY_OPTIONS: Array<{
	value: PostPrivacy;
	label: string;
	description: string;
}> = [
	{
		value: 'public',
		label: 'Public',
		description: 'Visible to everyone who can access your profile.',
	},
	{
		value: 'almost_private',
		label: 'Almost Private',
		description: 'Visible to your followers only.',
	},
	{
		value: 'private',
		label: 'Private',
		description: 'Visible only to followers you choose below.',
	},
];

export default function FeedPage() {
	const [feedData, setFeedData] = React.useState<PostListResponse | null>(null);
	const [loading, setLoading] = React.useState(true);
	const [feedError, setFeedError] = React.useState<string | null>(null);
	const [currentUser, setCurrentUser] = React.useState<PublicUser | null>(null);
	const [isComposerOpen, setIsComposerOpen] = React.useState(false);
	const [draftContent, setDraftContent] = React.useState('');
	const [postPrivacy, setPostPrivacy] = React.useState<PostPrivacy>('public');
	const [selectedViewerIds, setSelectedViewerIds] = React.useState<string[]>([]);
	const [selectedImage, setSelectedImage] = React.useState<File | null>(null);
	const [imagePreviewUrl, setImagePreviewUrl] = React.useState<string | null>(null);
	const [composerError, setComposerError] = React.useState<string | null>(null);
	const [isSubmitting, setIsSubmitting] = React.useState(false);
	const [followersData, setFollowersData] = React.useState<FollowListResponse | null>(null);
	const [followers, setFollowers] = React.useState<PublicUser[]>([]);
	const [followersLoading, setFollowersLoading] = React.useState(false);
	const [followersError, setFollowersError] = React.useState<string | null>(null);
	const [followersInitialized, setFollowersInitialized] = React.useState(false);
	const fileInputRef = React.useRef<HTMLInputElement | null>(null);

	const fetchFeed = async () => {
		setLoading(true);
		setFeedError(null);
		try {
			const data = await feedAPI.getFeed(DEFAULT_FEED_LIMIT, 0);
			setFeedData(data);
		} catch (err) {
			setFeedError(
				err instanceof Error ? err.message : 'Failed to load feed'
			);
		} finally {
			setLoading(false);
		}
	};

	React.useEffect(() => {
		const loadInitialData = async () => {
			setLoading(true);
			setFeedError(null);

			try {
				const [feed, profile] = await Promise.all([
					feedAPI.getFeed(DEFAULT_FEED_LIMIT, 0),
					profileAPI.getMyProfile(),
				]);
				setFeedData(feed);
				setCurrentUser(profile.user);
			} catch (err) {
				setFeedError(
					err instanceof Error ? err.message : 'Failed to load feed'
				);
			} finally {
				setLoading(false);
			}
		};

		void loadInitialData();
	}, []);

	React.useEffect(() => {
		if (!selectedImage) {
			setImagePreviewUrl(null);
			return;
		}

		const nextPreviewUrl = URL.createObjectURL(selectedImage);
		setImagePreviewUrl(nextPreviewUrl);

		return () => {
			URL.revokeObjectURL(nextPreviewUrl);
		};
	}, [selectedImage]);

	const loadFollowers = async (offset: number) => {
		if (!currentUser) {
			return;
		}

		setFollowersLoading(true);
		setFollowersError(null);

		try {
			const data = await profileAPI.getFollowers(
				currentUser.id,
				FOLLOWERS_PAGE_SIZE,
				offset
			);

			setFollowersData(data);
			setFollowers((current) => {
				if (offset === 0) {
					return data.users;
				}

				const nextUsers = [...current];
				for (const user of data.users) {
					if (!nextUsers.some((existingUser) => existingUser.id === user.id)) {
						nextUsers.push(user);
					}
				}
				return nextUsers;
			});
			setFollowersInitialized(true);
		} catch (err) {
			setFollowersError(
				err instanceof Error ? err.message : 'Failed to load followers'
			);
		} finally {
			setFollowersLoading(false);
		}
	};

	React.useEffect(() => {
		if (postPrivacy !== 'private' || !currentUser || followersInitialized) {
			return;
		}

		const loadInitialFollowers = async () => {
			setFollowersLoading(true);
			setFollowersError(null);

			try {
				const data = await profileAPI.getFollowers(
					currentUser.id,
					FOLLOWERS_PAGE_SIZE,
					0
				);
				setFollowersData(data);
				setFollowers(data.users);
				setFollowersInitialized(true);
			} catch (err) {
				setFollowersError(
					err instanceof Error ? err.message : 'Failed to load followers'
				);
			} finally {
				setFollowersLoading(false);
			}
		};

		void loadInitialFollowers();
	}, [postPrivacy, currentUser, followersInitialized]);

	const formatDate = (dateStr: string) => {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins} minute${diffMins === 1 ? '' : 's'} ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
	};

	const formatPrivacyLabel = (privacy: string) => {
		switch (privacy) {
			case 'almost_private':
				return 'almost private';
			default:
				return privacy;
		}
	};

	const resetComposer = () => {
		setDraftContent('');
		setPostPrivacy('public');
		setSelectedViewerIds([]);
		setSelectedImage(null);
		setImagePreviewUrl(null);
		setComposerError(null);
		setIsComposerOpen(false);
		if (fileInputRef.current) {
			fileInputRef.current.value = '';
		}
	};

	const openComposer = () => {
		setIsComposerOpen(true);
		setComposerError(null);
	};

	const openImagePicker = () => {
		openComposer();
		fileInputRef.current?.click();
	};

	const handleImageSelection = (event: React.ChangeEvent<HTMLInputElement>) => {
		const nextFile = event.target.files?.[0] ?? null;
		setSelectedImage(nextFile);
		if (nextFile) {
			setComposerError(null);
			setIsComposerOpen(true);
		}
	};

	const prependPost = (post: PostResponse) => {
		setFeedData((current) => {
			if (!current) {
				return {
					posts: [post],
					total: 1,
					limit: DEFAULT_FEED_LIMIT,
					offset: 0,
				};
			}

			const nextPosts = [post, ...current.posts];
			const cappedPosts = nextPosts.slice(0, current.limit);

			return {
				...current,
				posts: cappedPosts,
				total: current.total + 1,
			};
		});
	};

	const toggleViewerSelection = (viewerId: string) => {
		setSelectedViewerIds((current) =>
			current.includes(viewerId)
				? current.filter((id) => id !== viewerId)
				: [...current, viewerId]
		);
	};

	const handlePrivacyChange = (privacy: PostPrivacy) => {
		setPostPrivacy(privacy);
		setComposerError(null);
		if (privacy !== 'private') {
			setSelectedViewerIds([]);
		}
	};

	const handleCreatePost = async (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		const trimmedContent = draftContent.trim();
		if (!trimmedContent && !selectedImage) {
			setComposerError('Write something or choose an image before posting.');
			return;
		}

		if (postPrivacy === 'private' && selectedViewerIds.length === 0) {
			setComposerError('Choose at least one follower for a private post.');
			return;
		}

		setIsSubmitting(true);
		setComposerError(null);
		setFeedError(null);

		try {
			const createdPost = await feedAPI.createPost({
				content: trimmedContent || undefined,
				privacy: postPrivacy,
				visible_to: postPrivacy === 'private' ? selectedViewerIds : undefined,
			});

			let uploadWarning: string | null = null;

			if (selectedImage) {
				try {
					await feedAPI.uploadPostImage(createdPost.id, selectedImage);
				} catch (err) {
					uploadWarning =
						err instanceof Error
							? `Your post was created, but the image upload failed: ${err.message}`
							: 'Your post was created, but the image upload failed.';
				}
			}

			const completedPost = await feedAPI.getPost(createdPost.id);
			prependPost(completedPost);
			resetComposer();

			if (uploadWarning) {
				setFeedError(uploadWarning);
			}
		} catch (err) {
			setComposerError(
				err instanceof Error ? err.message : 'Failed to create post'
			);
		} finally {
			setIsSubmitting(false);
		}
	};

	const renderAvatar = (user?: PublicUser) => {
		const avatarUrl = resolveAssetUrl(user?.avatar_path);
		const initials = `${user?.first_name?.[0] ?? ''}${user?.last_name?.[0] ?? ''}`.trim() || 'Y';

		if (avatarUrl) {
			return (
				<img
					src={avatarUrl}
					alt={user ? `${user.first_name} ${user.last_name}` : 'User avatar'}
					className="w-10 h-10 rounded-full object-cover flex-shrink-0"
				/>
			);
		}

		return (
			<div className="w-10 h-10 rounded-full bg-indigo-100 text-indigo-700 text-sm font-semibold flex items-center justify-center flex-shrink-0">
				{initials}
			</div>
		);
	};

	const hasMoreFollowers =
		followersData !== null && followers.length < followersData.total;

	return (
		<div className="space-y-6">
			<input
				ref={fileInputRef}
				type="file"
				accept="image/*"
				className="hidden"
				onChange={handleImageSelection}
			/>

			<div className="bg-white p-4 rounded-xl border border-gray-100">
				{isComposerOpen ? (
					<form onSubmit={handleCreatePost} className="flex gap-4">
						{renderAvatar(currentUser ?? undefined)}
						<div className="flex-1 space-y-4">
							<textarea
								value={draftContent}
								onChange={(event) => setDraftContent(event.target.value)}
								placeholder="What&apos;s on your mind?"
								rows={4}
								className="w-full resize-none rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-800 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
								disabled={isSubmitting}
							/>

							<div className="rounded-xl border border-gray-100 bg-gray-50 p-3 space-y-3">
								<div>
									<p className="text-sm font-semibold text-gray-800">Post audience</p>
									<p className="text-xs text-gray-500">
										Choose who can see this post.
									</p>
								</div>
								<div className="grid gap-2 sm:grid-cols-3">
									{PRIVACY_OPTIONS.map((option) => {
										const isSelected = postPrivacy === option.value;

										return (
											<button
												key={option.value}
												type="button"
												onClick={() => handlePrivacyChange(option.value)}
												className={`rounded-xl border px-4 py-3 text-left transition-colors ${
													isSelected
														? 'border-indigo-200 bg-indigo-50'
														: 'border-gray-200 bg-white hover:bg-gray-50'
												}`}
												disabled={isSubmitting}
											>
												<div className="flex items-center justify-between gap-3 mb-1">
													<span className="text-sm font-semibold text-gray-800">
														{option.label}
													</span>
													{isSelected && <Check className="w-4 h-4 text-indigo-600" />}
												</div>
												<p className="text-xs text-gray-500 leading-relaxed">
													{option.description}
												</p>
											</button>
										);
									})}
								</div>
							</div>

							{postPrivacy === 'private' && (
								<div className="rounded-xl border border-gray-100 bg-gray-50 p-3 space-y-3">
									<div className="flex items-center justify-between gap-3">
										<div>
											<p className="text-sm font-semibold text-gray-800">
												Choose viewers
											</p>
											<p className="text-xs text-gray-500">
												Select which followers can see this post.
											</p>
										</div>
										<span className="text-xs font-medium text-indigo-700 bg-indigo-100 rounded-full px-3 py-1">
											{selectedViewerIds.length} selected
										</span>
									</div>

									{followersError && (
										<div className="rounded-xl border border-red-200 bg-white px-4 py-3">
											<p className="text-sm text-red-600">{followersError}</p>
										</div>
									)}

									{followersLoading && followers.length === 0 ? (
										<div className="flex items-center justify-center rounded-xl border border-gray-200 bg-white py-8">
											<LoaderCircle className="w-5 h-5 text-indigo-600 animate-spin" />
										</div>
									) : followers.length === 0 ? (
										<div className="rounded-xl border border-gray-200 bg-white px-4 py-6 text-center">
											<p className="text-sm text-gray-500">
												You do not have any followers to choose from yet.
											</p>
										</div>
									) : (
										<>
											<div className="max-h-64 overflow-y-auto rounded-xl border border-gray-200 bg-white divide-y divide-gray-100">
												{followers.map((follower) => {
													const isSelected = selectedViewerIds.includes(follower.id);
													const displayName = `${follower.first_name} ${follower.last_name}`;

													return (
														<button
															key={follower.id}
															type="button"
															onClick={() => toggleViewerSelection(follower.id)}
															className="w-full px-4 py-3 flex items-center gap-3 text-left hover:bg-gray-50 transition-colors"
															disabled={isSubmitting}
														>
															{renderAvatar(follower)}
															<div className="min-w-0 flex-1">
																<p className="text-sm font-semibold text-gray-800 truncate">
																	{displayName}
																</p>
																<p className="text-xs text-gray-500 truncate">
																	{follower.nickname ? `@${follower.nickname}` : 'Follower'}
																</p>
															</div>
															<div
																className={`w-5 h-5 rounded border flex items-center justify-center ${
																	isSelected
																		? 'bg-indigo-600 border-indigo-600'
																		: 'bg-white border-gray-300'
																}`}
															>
																{isSelected && <Check className="w-3 h-3 text-white" />}
															</div>
														</button>
													);
												})}
											</div>

											{hasMoreFollowers && (
												<div className="flex justify-center">
													<Button
														type="button"
														variant="secondary"
														size="sm"
														className="rounded-full"
														onClick={() => void loadFollowers(followers.length)}
														disabled={followersLoading || isSubmitting}
													>
														{followersLoading ? (
															<>
																<LoaderCircle className="w-4 h-4 mr-2 animate-spin" />
																Loading...
															</>
														) : (
															<>
																<ChevronDown className="w-4 h-4 mr-2" />
																Load More
															</>
														)}
													</Button>
												</div>
											)}
										</>
									)}
								</div>
							)}

							{imagePreviewUrl && (
								<div className="rounded-xl border border-gray-100 bg-gray-50 p-3">
									<div className="flex items-start justify-between gap-3 mb-3">
										<div>
											<p className="text-sm font-semibold text-gray-800">
												Selected image
											</p>
											<p className="text-xs text-gray-500">
												{selectedImage?.name}
											</p>
										</div>
										<Button
											type="button"
											variant="ghost"
											size="sm"
											className="rounded-full w-9 p-0 text-gray-500"
											onClick={() => setSelectedImage(null)}
											disabled={isSubmitting}
										>
											<X className="w-4 h-4" />
										</Button>
									</div>
									<img
										src={imagePreviewUrl}
										alt="Selected post preview"
										className="w-full max-h-80 rounded-xl object-cover border border-gray-100"
									/>
								</div>
							)}

							{composerError && (
								<div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3">
									<p className="text-sm text-red-600">{composerError}</p>
								</div>
							)}

							<div className="flex flex-wrap items-center justify-between gap-3">
								<Button
									type="button"
									variant="ghost"
									size="sm"
									className="rounded-full px-4 text-gray-600"
									onClick={openImagePicker}
									disabled={isSubmitting}
								>
									<ImageIcon className="w-4 h-4 mr-2 text-gray-400" />
									Add Image
								</Button>

								<div className="flex items-center gap-2">
									<Button
										type="button"
										variant="secondary"
										size="sm"
										className="rounded-full"
										onClick={resetComposer}
										disabled={isSubmitting}
									>
										Cancel
									</Button>
									<Button
										type="submit"
										variant="primary"
										size="sm"
										className="rounded-full px-5"
										disabled={isSubmitting}
									>
										{isSubmitting ? (
											<>
												<LoaderCircle className="w-4 h-4 mr-2 animate-spin" />
												Posting...
											</>
										) : (
											'Post'
										)}
									</Button>
								</div>
							</div>
						</div>
					</form>
				) : (
					<div className="flex gap-4">
						{renderAvatar(currentUser ?? undefined)}
						<button
							type="button"
							onClick={openComposer}
							className="flex-1 bg-gray-50 rounded-full px-5 text-left text-gray-500 text-sm hover:bg-gray-100 transition-colors"
						>
							What&apos;s on your mind?
						</button>
						<div className="flex gap-2">
							<Button
								type="button"
								variant="ghost"
								size="sm"
								className="rounded-full w-10 p-0"
								onClick={openImagePicker}
							>
								<ImageIcon className="w-5 h-5 text-gray-400" />
							</Button>
						</div>
					</div>
				)}
			</div>

			{feedError && (
				<div className="bg-white rounded-xl border border-red-200 p-4 text-center">
					<p className="text-sm text-red-600">{feedError}</p>
				</div>
			)}

			{loading && !feedData && (
				<div className="space-y-4">
					{[1, 2, 3].map((i) => (
						<div key={i} className="bg-white rounded-xl border border-gray-100 p-4 space-y-3">
							<div className="flex items-center gap-3">
								<div className="w-10 h-10 rounded-full bg-gray-200 animate-pulse" />
								<div className="space-y-1">
									<div className="w-32 h-4 bg-gray-200 rounded animate-pulse" />
									<div className="w-20 h-3 bg-gray-100 rounded animate-pulse" />
								</div>
							</div>
							<div className="space-y-2">
								<div className="w-full h-4 bg-gray-100 rounded animate-pulse" />
								<div className="w-3/4 h-4 bg-gray-100 rounded animate-pulse" />
							</div>
						</div>
					))}
				</div>
			)}

			{!loading && feedData && feedData.posts.length === 0 && (
				<div className="bg-white rounded-xl border border-gray-100 p-8 text-center">
					<FileText className="w-12 h-12 text-gray-300 mx-auto mb-3" />
					<h3 className="text-lg font-semibold text-gray-800 mb-1">
						No posts to display
					</h3>
					<p className="text-sm text-gray-500">
						Your feed is empty. Follow users to see their posts here.
					</p>
				</div>
			)}

			{feedData && feedData.posts.length > 0 && (
				<>
					{feedData.posts.map((post) => {
						const postImageUrl = resolveAssetUrl(post.image_url);

						return (
							<motion.div
								key={post.id}
								initial={{ opacity: 0, y: 10 }}
								animate={{ opacity: 1, y: 0 }}
								className="bg-white rounded-xl border border-gray-100 overflow-hidden"
							>
								<div className="p-4 flex items-center justify-between">
									<div className="flex items-center gap-3">
										{renderAvatar(post.author)}
										<div>
											<p className="text-sm font-bold">
												{post.author.first_name} {post.author.last_name}
											</p>
											<p className="text-xs text-gray-500">
												{formatDate(post.created_at)} • {formatPrivacyLabel(post.privacy)}
											</p>
										</div>
									</div>
									<button className="text-gray-400 hover:text-gray-600">
										<MoreHorizontal className="w-5 h-5" />
									</button>
								</div>

								{post.content && (
									<div className="px-4 pb-4">
										<p className="text-sm leading-relaxed text-gray-800">
											{post.content}
										</p>
									</div>
								)}

								{postImageUrl && (
									<div className="px-4 pb-4">
										<img
											src={postImageUrl}
											alt="Post image"
											className="w-full max-h-[32rem] rounded-xl border border-gray-100 object-cover bg-gray-50"
										/>
									</div>
								)}

								<div className="px-4 py-3 bg-gray-50/50 border-t border-gray-100 flex items-center gap-6">
									<button className="flex items-center gap-2 text-xs font-bold text-gray-500 hover:text-indigo-600 transition-colors">
										<MessageCircle className="w-4 h-4" />
										Comment
									</button>
									<button className="flex items-center gap-2 text-xs font-bold text-gray-500 hover:text-indigo-600 transition-colors">
										<Share2 className="w-4 h-4" />
										Share
									</button>
								</div>
							</motion.div>
						);
					})}
				</>
			)}
		</div>
	);
}
