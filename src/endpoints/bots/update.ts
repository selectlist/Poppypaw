import * as database from "../../Serendipity/prisma.js";
import { FastifyReply, FastifyRequest } from "fastify";

export default {
	method: "PATCH",
	url: "/bots/update",
	schema: {
		summary: "Update Bot",
		description:
			"Returns boolean value indicating whether the update was successful or not.",
		tags: ["bots"],
		body: {
			type: "object",
			properties: {
				platform: { type: "string" },
				bot_id: { type: "string" },
				description: { type: "string" },
				long_description: { type: "string" },
			},
			required: ["platform", "bot_id", "description", "long_description"],
		},
		security: [
			{
				apiKey: [],
			},
		],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		try {
			const { platform, bot_id, description, long_description }: any =
				request.body;
			const Authorization: any = request.headers.authorization;
			const token = await database.Tokens.get({ token: Authorization });
			const dbUser = await database.Users.get({ userid: token.userid });

			if (dbUser) {
				if (platform === "discord") {
					let bot = await database.Discord.get({
						botid: bot_id,
					});

					if (bot.ownerid === dbUser.userid) {
						bot.description = description;
						bot.longdescription = long_description;

						const update = await database.Discord.update(
							bot_id,
							bot
						);

						if (update) return reply.status(204).send();
						else
							return reply.status(500).send({
								error: "An error has occured while proccessing your request. This error has automatically been reported to our team.",
							});
					} else
						return reply.send({
							token: Authorization,
							error: true,
							message: "You do not have access to this bot.",
						});
				} else if (platform === "revolt") {
					const bot = await database.Revolt.get({
						botid: bot_id,
					});

					if (bot.ownerid === dbUser.revoltid) {
						bot.description = description;
						bot.longdescription = long_description;

						const update = await database.Revolt.update(
							bot_id,
							bot
						);

						if (update) return reply.status(204).send();
						else
							return reply.status(500).send({
								error: "An error has occured while proccessing your request. This error has automatically been reported to our team.",
							});
					} else
						return reply.send({
							token: Authorization,
							error: true,
							message: "You do not have access to this bot.",
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
