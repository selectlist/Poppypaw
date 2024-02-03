import { FastifyReply, FastifyRequest } from "fastify";
import discord from "discord-oauth2";
import * as crypto from "crypto";

export default {
	url: "/auth/login",
	method: "GET",
	schema: {
		summary: "Login",
		description:
			"This endpoint will redirect you to Discord oAuth. The callback will automatically be determined by the Request Referrer. Please note that this will only work if your domain is whitelisted.",
		tags: ["auth"],
		querystring: {
			type: "object",
			properties: {
				redirect: { type: "string" },
			},
			required: ["redirect"],
		},
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		const oauth = new discord();
		const { redirect }: any = request.query;

		const state = JSON.stringify({
			csrf_token: crypto
				.createHash("sha256")
				.update(
					`${crypto.randomUUID()}_${crypto.randomUUID()}`.replace(
						/-/g,
						""
					)
				)
				.digest("hex"),
			date: new Date().toUTCString(),
			ip: request.ip,
			user_agent: request.headers["user-agent"],
			redirect: redirect,
		});

		const url = oauth.generateAuthUrl({
			clientId: process.env.DISCORD_CLIENT_ID,
			redirectUri: process.env.DISCORD_AUTH_REDIRECT,
			scope: ["identify", "guilds"],
			state: encodeURIComponent(state),
		});

		reply.redirect(url);
	},
};
