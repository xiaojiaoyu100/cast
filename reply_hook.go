package cast

type ReplyHook func(cast *Cast, reply *Reply) error

func dumpReply(cast *Cast, reply *Reply) error {
	if cast == nil || reply == nil {
		return nil
	}

	if cast.logger == nil {
		return nil
	}

	cast.logger.Printf("%s took %s upto %d time(s)", reply.Url(), reply.cost, reply.times)

	return nil
}
