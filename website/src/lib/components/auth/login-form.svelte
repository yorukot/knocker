<script lang="ts">
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { buildOAuthUrl, login } from '$lib/api/auth.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import {
		FieldGroup,
		Field,
		FieldLabel,
		FieldDescription
	} from '$lib/components/ui/field/index.js';
	import { z } from 'zod';

	const inBrowser = typeof window !== 'undefined';

	let redirectTo = $state('/');
	const form = $state({
		email: '',
		password: '',
		loading: false,
		error: '',
		success: ''
	});

	const loginSchema = z.object({
		email: z.email().max(255),
		password: z.string().min(8).max(255)
	});

	if (inBrowser) {
		const url = new URL(window.location.href);
		redirectTo = url.searchParams.get('next') ?? '/';
	}

	const handleSubmit = async (event: Event) => {
		event.preventDefault();
		form.error = '';
		form.success = '';

		const parsed = loginSchema.safeParse({
			email: form.email.trim(),
			password: form.password
		});

		if (!parsed.success) {
			form.error = parsed.error.issues[0]?.message ?? 'Please check your inputs.';
			return;
		}

		form.loading = true;

		try {
			await login(parsed.data.email, parsed.data.password);
			form.success = 'Logged in successfully.';

			if (inBrowser) {
				await goto(redirectTo);
			}
		} catch (error) {
			form.error =
				error instanceof Error ? error.message : 'Unable to log in right now. Please try again.';
		} finally {
			form.loading = false;
		}
	};

	const handleGoogle = () => {
		if (!inBrowser) return;
		window.location.href = buildOAuthUrl('google', redirectTo);
	};
</script>

<Card.Root class="mx-auto w-full max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Login</Card.Title>
		<Card.Description>Enter your email below to login to your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<form class="space-y-4" onsubmit={handleSubmit}>
			<FieldGroup>
				<Field>
					<FieldLabel for="email">Email</FieldLabel>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="m@example.com"
						autocomplete="email"
						bind:value={form.email}
						required
					/>
				</Field>
				<Field>
					<div class="flex items-center">
						<FieldLabel for="password">Password</FieldLabel>
						<a href="##" class="ms-auto inline-block text-sm underline"> Forgot your password? </a>
					</div>
					<Input
						id="password"
						name="password"
						type="password"
						autocomplete="current-password"
						bind:value={form.password}
						required
					/>
				</Field>
				{#if form.error}
					<FieldDescription class="text-destructive" data-invalid="true">
						{form.error}
					</FieldDescription>
				{:else if form.success}
					<FieldDescription class="text-emerald-600 dark:text-emerald-400">
						{form.success}
					</FieldDescription>
				{/if}
				<Field>
					<Button type="submit" class="w-full" disabled={form.loading}>
						{form.loading ? 'Logging in...' : 'Login'}
					</Button>
					<Button type="button" variant="outline" class="w-full" onclick={handleGoogle}>
						<Icon icon="ri:google-fill" class="size-5" />
						Login with Google
					</Button>
					<FieldDescription class="text-center">
						Don't have an account? <a href="/auth/register">Sign up</a>
					</FieldDescription>
				</Field>
			</FieldGroup>
		</form>
	</Card.Content>
</Card.Root>
