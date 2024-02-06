import * as database from "../../Serendipity/prisma.js";
import { FastifyReply, FastifyRequest } from "fastify";
import { REST } from "@discordjs/rest";
import { Routes } from "discord-api-types/v9";
import type { RESTGetAPIUserResult } from "discord-api-types/v9";
import { Client } from "revolt.js";
import { info } from "../../logger.js";

// Initalize REST
const rest = new REST({
	version: "9",
}).setToken(process.env.DISCORD_TOKEN as string);

// Create a Revolt Chat client
const revoltClient: Client = new Client();
revoltClient.loginBot(process.env.REVOLT_TOKEN);
revoltClient.on("ready", async () =>
	info("Revolt (Dovewing - Bot Add)", "Ready!")
);

export default {
	method: "POST",
	url: "/bots/add",
	schema: {
		summary: "Add Bot",
		description:
			"Returns boolean value indicating whether the creation was successful or not.",
		tags: ["bots"],
		body: {
			type: "object",
			properties: {
				platform: { type: "string" },
				bot_id: { type: "string" },
				invite: { type: "string" },
				description: { type: "string" },
				long_description: { type: "string" },
			},
			required: [
				"platform",
				"bot_id",
				"invite",
				"description",
				"long_description",
			],
		},
		security: [
			{
				apiKey: [],
			},
		],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		try {
			const {
				platform,
				bot_id,
				invite,
				description,
				long_description,
			}: any = request.body;
			const Authorization: any = request.headers.authorization;
			const token = await database.Tokens.get({ token: Authorization });
			const dbUser = await database.Users.get({ userid: token.userid });

			if (dbUser) {
				if (platform === "discord") {
					const apiUserData = (await rest.get(
						Routes.user(bot_id)
					)) as RESTGetAPIUserResult;

					if (!apiUserData)
						return reply.status(400).send({
							error: "Bot does not exist.",
							status: 400,
						});
					else if (apiUserData.bot) {
						const data = await database.Discord.create({
							botid: bot_id,
							name: apiUserData.username,
							avatar: `https://cdn.discordapp.com/avatars/${bot_id}/${apiUserData.avatar}.png`,
							invite: invite,
							description: description,
							longdescription: long_description,
							servers: 0,
							shards: 0,
							users: 0,
							claimedBy: null,
							state: "PENDING",
							upvotes: [],
							downvotes: [],
							ownerid: dbUser.userid,
						});

						if (data) return reply.status(204).send();
						else
							return reply.status(500).send({
								error: "An error has occured while proccessing your request. This error has automatically been reported to our team.",
							});
					}
				} else if (platform === "revolt") {
					await revoltClient.users.fetch(bot_id).then(async (p) => {
						if (!p)
							return reply.status(400).send({
								error: "Bot does not exist.",
								status: 400,
							});
						else if (p.bot) {
							const data = await database.Revolt.create({
								botid: bot_id,
								name: p.username,
								avatar: p.avatarURL,
								invite: invite,
								description: description,
								longdescription: long_description,
								servers: 0,
								shards: 0,
								users: 0,
								claimedBy: null,
								state: "PENDING",
								upvotes: [],
								downvotes: [],
								ownerid: dbUser.revoltid,
							});

							if (data) return reply.status(204).send();
							else
								return reply.status(500).send({
									error: "An error has occured while proccessing your request. This error has automatically been reported to our team.",
								});
						}
					});
				}
			} else
				return reply.send({
					token: Authorization,
					error: true,
					message: "User does not exist.",
				});
		} catch (error) {
			reply.status(500).send({
				error: "Internal Server Error",
				message: error.errorInfo.message,
			});
		}
	},
};
