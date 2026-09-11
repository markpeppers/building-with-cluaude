"""
Algorithm to calculate Pi to 5 decimal places using the Leibniz formula.
The Leibniz formula: π/4 = 1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
"""

def calculate_pi(precision=5):
    """
    Calculate pi to the specified number of decimal places using the Leibniz formula.
    
    Args:
        precision: Number of decimal places (default: 5)
    
    Returns:
        float: Approximation of pi
    """
    # We need enough iterations to get the desired precision
    # For 5 decimal places, we need more iterations
    max_iterations = 500000
    
    pi_over_4 = 0.0
    
    for i in range(max_iterations):
        # Leibniz formula: π/4 = 1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
        term = ((-1) ** i) / (2 * i + 1)
        pi_over_4 += term
    
    pi_approximation = pi_over_4 * 4
    
    # Round to the specified precision
    return round(pi_approximation, precision)


def main():
    """Main function to calculate and display pi."""
    pi_value = calculate_pi(5)
    print(f"Pi calculated to 5 decimal places: {pi_value}")
    print(f"Pi (for reference): 3.14159")
    return pi_value


if __name__ == "__main__":
    main()
