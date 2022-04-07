from sym.sdk.strategies import AccessStrategy


class CustomAccessStrategy(AccessStrategy):
    """A completely no-op strategy, but technically valid.
    """

    def fetch_remote_identity(self, user):
        return None

    def escalate(self, target_id, event):
        pass

    def deescalate(self, target_id, event):
        pass
