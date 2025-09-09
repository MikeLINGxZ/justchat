
export class WaitGroup {
    private counter = 0;
    private resolve!: () => void;
    private promise: Promise<void> = new Promise((res) => (this.resolve = res));

    add(delta: number): void {
        this.counter += delta;
        if (this.counter < 0) {
            throw new Error("WaitGroup counter cannot be negative");
        }
        if (this.counter === 0) {
            this.resolve();
        }
    }

    done(): void {
        this.add(-1);
    }

    wait(): Promise<void> {
        return this.promise;
    }
}