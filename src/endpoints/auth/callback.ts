import { FastifyReply, FastifyRequest } from "fastify";
import discord from "discord-oauth2";
import * as crypto from "crypto";
import * as database from "../../Serendipity/prisma.js";

export default {
	url: "/auth/callback",
	method: "GET",
	schema: {
		summary: "Auth Callback",
		description:
			"Authorization Callback Endpoint for Discord oAuth, based off of the /auth/login endpoint.",
		tags: ["auth"],
		querystring: {
			type: "object",
			properties: {
				code: { type: "string" },
				state: { type: "string" },
			},
			required: ["code", "state"],
		},
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		const oauth = new discord();

		const { code, state }: any = request.query;
		if (!code || !state)
			return reply.send({
				error: "No code/state provided from Discord.",
			});

		const { redirect, user_agent } = JSON.parse(
			decodeURIComponent(state as string)
		);

		const token = await oauth.tokenRequest({
			clientId: process.env.DISCORD_CLIENT_ID,
			clientSecret: process.env.DISCORD_CLIENT_SECRET,

			code: code,
			scope: "identify guilds",
			grantType: "authorization_code",

			redirectUri: process.env.DISCORD_AUTH_REDIRECT,
		});

		const discordUser = await oauth.getUser(token.access_token);
		if (!discordUser)
			return reply.redirect(`/auth/login?redirect=${redirect}`);

		const authCode = crypto
			.createHash("sha256")
			.update(
				`${crypto.randomUUID()}_${crypto.randomUUID()}_${crypto.randomUUID()}_${crypto.randomUUID()}`.replace(
					/-/g,
					""
				)
			)
			.digest("hex");

		const user = await database.Users.get({
			userid: discordUser.id,
		});

		if (user) {
			await database.Tokens.create({
				id: crypto.randomUUID().toString(),
				userid: discordUser.id,
				token: authCode,
				agent: user_agent,
				createdAt: new Date(),
			});
		} else {
			await database.Users.create({
				username: discordUser.username,
				userid: discordUser.id,
				revoltid: null,
				bio: null,
				avatar: `https://cdn.discordapp.com/avatars/${discordUser.id}/${discordUser.avatar}.png`,
				badges: [],
				staff_perms: [],
			});

			await database.Tokens.create({
				id: crypto.randomUUID().toString(),
				userid: discordUser.id,
				token: authCode,
				agent: user_agent,
				createdAt: new Date(),
			});
		}

		return reply.redirect(`${redirect}?auth_code=${authCode}`);
	},
};
