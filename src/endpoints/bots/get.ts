import { FastifyReply, FastifyRequest } from "fastify";
import * as database from "../../v4-database/prisma.js";

export default {
	url: "/bots/get",
	method: "GET",
	schema: {
		summary: "Get Bot",
		description: "Gets a bot.",
		tags: ["bots"],
		querystring: {
			type: "object",
			properties: {
				botid: { type: "string" },
			},
			required: ["botid"],
		},
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		const data: any = request.query;

		let bot = await database.Discord.get({ botid: data.botid });

		if (bot) return reply.send(bot);
		else {
			bot = await database.Revolt.get({ botid: data.botid });

			if (bot) return reply.send(bot);
			else
				return reply.status(404).send({
					message:
						"We couldn't fetch any information about this bot in our database",
					error: true,
				});
		}
	},
};
