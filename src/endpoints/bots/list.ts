import * as database from "../../v4-database/prisma.js";
import { FastifyReply, FastifyRequest } from "fastify";

export default {
	url: "/bots/list",
	method: "GET",
	schema: {
		summary: "Get all bots",
		description: "Returns all bots.",
		tags: ["bots"],
	},
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		let discord = await database.Discord.find({});
		discord.reverse();

		let revolt = await database.Revolt.find({});
		revolt.reverse();

		return reply.send({
			discord: discord,
			revolt: revolt,
		});
	},
};
