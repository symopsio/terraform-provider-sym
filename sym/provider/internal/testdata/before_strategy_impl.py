from sym.sdk.strategies import AccessStrategy


class CustomAccessStrategy(AccessStrategy):
    """A completely no-op strategy, but technically valid.

    This is some extra commenting so the file contents are technically
    different than the before_strategy_impl.py contents.
    """

    def fetch_remote_identity(self, user):
        return None

    def escalate(self, target_id, event):
        pass

    def deescalate(self, target_id, event):
        pass
