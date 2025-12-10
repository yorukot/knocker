<script lang="ts">
	import { goto } from '$app/navigation';
	import { buildOAuthUrl, registerUser } from '$lib/api/auth.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Icon from '@iconify/svelte';
	import type { ComponentProps } from 'svelte';
	import { z } from 'zod';

	let { ...restProps }: ComponentProps<typeof Card.Root> = $props();
	const inBrowser = typeof window !== 'undefined';
	let redirectTo = $state('/');
	const form = $state({
		displayName: '',
		email: '',
		password: '',
		confirmPassword: '',
		loading: false,
		error: '',
		success: ''
	});

	const signupSchema = z
		.object({
			displayName: z.string().min(3).max(255),
			email: z.email().max(255),
			password: z.string().min(8).max(255),
			confirmPassword: z.string().min(8).max(255)
		})
		.refine((value) => value.password === value.confirmPassword, {
			path: ['confirmPassword'],
			message: 'Passwords do not match.'
		});

	if (inBrowser) {
		const url = new URL(window.location.href);
		redirectTo = url.searchParams.get('next') ?? '/';
	}

	const handleSubmit = async (event: Event) => {
		event.preventDefault();
		form.error = '';
		form.success = '';

		const parsed = signupSchema.safeParse({
			displayName: form.displayName.trim(),
			email: form.email.trim(),
			password: form.password,
			confirmPassword: form.confirmPassword
		});

		if (!parsed.success) {
			form.error = parsed.error.issues[0]?.message ?? 'Please check your inputs.';
			return;
		}

		form.loading = true;

		try {
			await registerUser(parsed.data.displayName, parsed.data.email, parsed.data.password);
			form.success = 'Account created. Redirecting...';

			if (inBrowser) {
				await goto(redirectTo);
			}
		} catch (error) {
			form.error =
				error instanceof Error ? error.message : 'Unable to sign up right now. Please try again.';
		} finally {
			form.loading = false;
		}
	};

	const handleGoogle = () => {
		if (!inBrowser) return;
		window.location.href = buildOAuthUrl('google', redirectTo);
	};
</script>

<Card.Root {...restProps}>
	<Card.Header>
		<Card.Title>Create an account</Card.Title>
		<Card.Description>Enter your information below to create your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<form class="space-y-4" onsubmit={handleSubmit}>
			<Field.Group>
				<Field.Field>
					<Field.Label for="name">Display Name</Field.Label>
					<Input
						id="name"
						name="display_name"
						type="text"
						placeholder="John Doe"
						autocomplete="name"
						bind:value={form.displayName}
						required
					/>
				</Field.Field>
				<Field.Field>
					<Field.Label for="email">Email</Field.Label>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="m@example.com"
						autocomplete="email"
						bind:value={form.email}
						required
					/>
					<Field.Description>
						We'll use this to contact you. We will not share your email with anyone else.
					</Field.Description>
				</Field.Field>
				<Field.Field>
					<Field.Label for="password">Password</Field.Label>
					<Input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						bind:value={form.password}
						required
					/>
					<Field.Description>Must be at least 8 characters long.</Field.Description>
				</Field.Field>
				<Field.Field>
					<Field.Label for="confirm-password">Confirm Password</Field.Label>
					<Input
						id="confirm-password"
						name="confirm_password"
						type="password"
						autocomplete="new-password"
						bind:value={form.confirmPassword}
						required
					/>
					<Field.Description>Please confirm your password.</Field.Description>
				</Field.Field>
				{#if form.error}
					<Field.Description class="text-destructive" data-invalid="true">
						{form.error}
					</Field.Description>
				{:else if form.success}
					<Field.Description class="text-emerald-600 dark:text-emerald-400">
						{form.success}
					</Field.Description>
				{/if}
				<Field.Group>
					<Field.Field>
						<Button type="submit" disabled={form.loading}>
							{form.loading ? 'Creating account...' : 'Create Account'}
						</Button>
						<Button variant="outline" type="button" onclick={handleGoogle}>
							<Icon icon="ri:google-fill" class="size-5" />
							Sign up with Google</Button
						>
						<Field.Description class="px-6 text-center">
							Already have an account? <a href="/auth/login">Sign in</a>
						</Field.Description>
					</Field.Field>
				</Field.Group>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
